package protogen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/internal/printer"
	"github.com/Wenchy/tableau/proto/tableaupb"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/proto"
	"github.com/pkg/errors"
)

type XmlGenerator struct {
	Generator

	fieldMap map[string]*tableaupb.Field
	nav *xmlquery.NodeNavigator
}

var numRegex *regexp.Regexp

func init() {
	numRegex = regexp.MustCompile(`[0-9]+`) // e.g.: Item1ID
}

func (gen *XmlGenerator) Generate() error {
	err := os.RemoveAll(gen.OutputDir)
	if err != nil {
		panic(err)
	}
	// create output dir
	err = os.MkdirAll(gen.OutputDir, 0700)
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(gen.InputDir)
	if err != nil {
		atom.Log.Fatal(err)
	}
	for _, xmlFile := range files {
		// ignore temp file named with prefix "~$"
		if strings.HasPrefix(xmlFile.Name(), "~$") {
			continue
		}
		// open xml file and parse the document
		xmlPath := filepath.Join(gen.InputDir, xmlFile.Name())
		atom.Log.Debugf("xml: %s", xmlPath)
		f, err := os.Open(xmlPath)
		if err != nil {
			atom.Log.Panic(err)
			continue
		}
		p, err := xmlquery.CreateStreamParser(f, "/")
		if err != nil {
			atom.Log.Panic(err)
			continue
		}
		// create xml proto meta struct
		xmlProtoName := strcase.ToSnake(strings.TrimSuffix(xmlFile.Name(), filepath.Ext(xmlFile.Name())))
		xml := &tableaupb.Workbook{
			Options: &tableaupb.WorkbookOptions{
				Name: xmlFile.Name(),
			},
			Name:       xmlProtoName,
			Imports: map[string]int32{
				tableauProtoPath: 1, // default import
			},
		}
		for _, path := range gen.Imports {
			xml.Imports[path] = 1 // custom imports
		}
		n, err := p.Read()
		if err != nil {
			atom.Log.Panic(err)
		}
		gen.fieldMap = make(map[string]*tableaupb.Field)
		worksheet := &tableaupb.Worksheet{
			Options: &tableaupb.WorksheetOptions{
				Name: xmlFile.Name(),
			},
			Name: xmlFile.Name(),
		}
		root := &tableaupb.Field{}
		gen.parseNode(xmlquery.CreateXPathNavigator(n), root, "")
		worksheet.Fields = append(worksheet.Fields, root)
		xml.Worksheets = append(xml.Worksheets, worksheet)		
		// export book
		be := newBookExporter(gen.ProtoPackage, gen.GoPackage, gen.OutputDir, gen.FilenameSuffix, gen.Imports, xml)
		if err := be.export(); err != nil {
			return errors.Wrapf(err, "failed to export workbook: %s", xmlPath)
		}
	}

	return nil
}

func (gen *XmlGenerator) parseNode(nav *xmlquery.NodeNavigator, element *tableaupb.Field, prefix string) error {	
	gen.nav = nav
	element.Options = &tableaupb.FieldOptions{
		Name: nav.LocalName(),
	}
	element.Name = nav.LocalName()
	// iterate over attributes
	for _, attr := range nav.Current().Attr {
		switch strings.ToLower(attr.Name.Local) {
		case "keycol":
			tagName := strings.Split(attr.Value, ".")[0]
			attrName := strings.Split(attr.Value, ".")[1]
			keyNode := xmlquery.FindOne(nav.Current(), fmt.Sprintf("%s/@%s", tagName, attrName))
			if keyNode == nil {
				atom.Log.Panic(fmt.Sprintf("KeyCol:%s not found in the immediately following nodes of %s", attr.Value, nav.LocalName()))
				continue
			}
			keyType, _ := gen.guessType(keyNode.InnerText())
			element.Options.Key = attr.Value
			element.MapEntry = &tableaupb.MapEntry{
				KeyType: keyType,
				ValueType: tagName,
			}
		case "desc":
		default:
			attrName := attr.Name.Local
			attrValue := attr.Value
			t, d := gen.guessType(attrValue)
			newAttr := &tableaupb.Field{
				Options: &tableaupb.FieldOptions{
					Name: attrName,
					Default: d,
				},
				Type: t,
				Name: strcase.ToSnake(attrName),
			}
			if matches := numRegex.FindStringSubmatch(attrName); len(matches) > 0 {
				if matches[0] != "1" {
					break
				}
				newAttr.Card = "repeated"
				newAttr.Options.Name = strings.ReplaceAll(newAttr.Options.Name, matches[0], "")
				newAttr.Name = strcase.ToSnake(newAttr.Options.Name)
			}
			element.Fields = append(element.Fields, newAttr)
		}
	}
	// fmt.Println(element)

	// createMetaStruct create a meta struct by some pre-defined tag
	tagMap := make(map[string]bool)
	createMetaStruct := func(metaNode *xmlquery.Node) {
		metaNav := xmlquery.CreateXPathNavigator(metaNode.FirstChild)
		for metaNav.NodeType() != xpath.ElementNode {
			metaNav.MoveToNext()
		}
		newChild := &tableaupb.Field{
			Options: &tableaupb.FieldOptions{
				Name: metaNav.LocalName(),
			},
			Card: "repeated",
			Type: metaNav.LocalName(),
			Name: strcase.ToSnake(metaNav.LocalName()) + "_list",
		}
		element.Fields = append(element.Fields, newChild)
		gen.parseNode(metaNav, newChild, fmt.Sprintf("%s/%s", prefix, metaNode.Parent.Data))
		gen.fieldMap[fmt.Sprintf("%s/%s/%s", prefix, metaNode.Parent.Data, metaNav.LocalName())] = newChild
		tagMap[metaNode.Data] = true
	}
	// `StructSupplement` defines the default values of one tag
	if metaNode := xmlquery.FindOne(nav.Current(), "StructSupplement"); metaNode != nil {		
		createMetaStruct(metaNode)
	}
	// `StructFormatSupplement` defines a meta struct
	if metaNode := xmlquery.FindOne(nav.Current(), "StructFormatSupplement"); metaNode != nil {		
		createMetaStruct(metaNode)
	}
	// iterate over child nodes
	navCopy := *nav
	for flag := navCopy.MoveToChild(); flag; flag = navCopy.MoveToNext() {
		// commentNode, documentNode and other meaningless nodes should be filtered
		if navCopy.NodeType() != xpath.ElementNode {
			continue
		}
		tagName := navCopy.LocalName()
		if _, exist := tagMap[tagName]; exist {
			continue
		}
		newChild := &tableaupb.Field{
			Options: &tableaupb.FieldOptions{
				Name: tagName,
			},
			Type: tagName,
			Name: strcase.ToSnake(tagName),
		}
		if keyTag, childList := strings.Split(element.Options.Key, ".")[0], xmlquery.Find(nav.Current(), tagName); len(childList) > 1 || tagName == keyTag {			
			newChild.Card = "repeated"
			newChild.Name = newChild.Name + "_list"
		}
		fatherPath := fmt.Sprintf("%s/%s", prefix, nav.Current().Data)
		gen.parseNode(&navCopy, newChild, fatherPath)
		
		// overwrite previous meta struct if necessary
		curPath := fatherPath + "/" + tagName
		if child, exist := gen.fieldMap[curPath]; exist {
			if child.Card == "" && newChild.Card == "repeated" {
				child.Card = "repeated"
				child.Name = child.Name + "_list"
			}
			fieldMap := make(map[string]bool)
			for _, c := range child.Fields {
				fieldMap[c.Options.Name] = true
			}
			for _, c := range newChild.Fields {
				if _, exist := fieldMap[c.Options.Name]; !exist {
					child.Fields = append(child.Fields, c)
				}
			}
		} else {
			gen.fieldMap[curPath] = newChild
			element.Fields = append(element.Fields, newChild)
		}
	}
	return nil
}

func (gen *XmlGenerator) guessType(value string) (string, string) {
	var t, d string
	if _, err := strconv.Atoi(value); err == nil {
		t, d = "int32", "0"
	} else if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		t, d = "int64", "0"
	} else {
		t, d = "string", ""
	}
	if gen.nav.Current().Parent.Data == "StructSupplement" {
		d = value
	}
	return t, d
}