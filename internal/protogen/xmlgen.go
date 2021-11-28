package protogen

import (
	// "github.com/antchfx/xpath"
	// "bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// "io"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"

	// "google.golang.org/protobuf/proto"
	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"

	// "github.com/antchfx/xpath"
	"github.com/iancoleman/strcase"
)

type XmlGenerator struct {
	ProtoPackage string // protobuf package name.
	GoPackage    string // golang package name.
	InputDir     string // input dir of workbooks.
	OutputDir    string // output dir of generated protoconf files.
	Imports      []string // imported common proto file paths

	Xml *options.XmlOption // xml generation settings
}

type xml struct {
	xml *tableaupb.Xml
}

var numRegex *regexp.Regexp

func init() {
	numRegex = regexp.MustCompile(`[0-9]+`) // e.g.: Item1ID
}

func (gen *XmlGenerator) Generate() {
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
		xml := &xml{
			xml: &tableaupb.Xml{
				Options: &tableaupb.XmlOptions{
					Name: xmlFile.Name(),
				},
				Root: &tableaupb.Element{},
				Name:       xmlProtoName,
				Imports: map[string]int32{
					tableauProtoPath: 1, // default import
				},
			},
		}
		for _, path := range gen.Imports {
			xml.xml.Imports[path] = 1 // custom imports
		}
		n, err := p.Read()
		if err != nil {
			atom.Log.Panic(err)
		}
		gen.parseNode(xmlquery.CreateXPathNavigator(n), xml.xml.Root)
	}
}

func (gen *XmlGenerator) parseNode(nav *xmlquery.NodeNavigator, element *tableaupb.Element) error {	
	element.Options = &tableaupb.ElementOptions{
		Name: nav.LocalName(),
	}
	// iterate over attributes
	for _, attr := range nav.Current().Attr {
		switch strings.ToLower(attr.Name.Local) {
		case "keycol":
			tagName := strings.Split(attr.Value, ".")[0]
			attrName := strings.Split(attr.Value, ".")[1]
			if xmlquery.FindOne(nav.Current(), fmt.Sprintf("%s/@%s", tagName, attrName)) == nil {
				atom.Log.Panic(fmt.Sprintf("KeyCol:%s not found in the immediately following nodes of %s", attr.Value, nav.LocalName()))
			}
			element.Options.Key = attrName
		case "desc":
		default:
			attrName := attr.Name.Local
			// attrValue := attr.Value
			newAttr := &tableaupb.Attr{
				Options: &tableaupb.AttrOptions{
					Name: attrName,
					Default: "", //TODO
				},
				Type: "", //TODO
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
			element.Attrs = append(element.Attrs, newAttr)
		}
	}
	fmt.Println(element)
	// `StructSupplement` defines the meta struct of one tag
	tagMap := make(map[string]bool)
	if metaNode := xmlquery.FindOne(nav.Current(), "StructSupplement"); metaNode != nil {		
		metaNav := xmlquery.CreateXPathNavigator(metaNode.FirstChild)
		for metaNav.NodeType() != xpath.ElementNode {
			metaNav.MoveToNext()
		}
		newChild := &tableaupb.Child{
			Options: &tableaupb.ChildOptions{
				Name: metaNav.LocalName(),
			},
			Card: "repeated",
			Type: metaNav.LocalName(),
			Name: strcase.ToSnake(metaNav.LocalName()) + "_list",
			Element: &tableaupb.Element{},
		}
		element.Children = append(element.Children, newChild)
		gen.parseNode(metaNav, newChild.Element)
		tagMap = map[string]bool{
			"StructSupplement": true,
			metaNav.LocalName(): true,
		}
	}
	// iterate over child nodes
	navCopy := *nav
	for flag := navCopy.MoveToChild(); flag; flag = navCopy.MoveToNext() {
		// context node must be the leaf node
		if content := strings.TrimSpace(navCopy.LocalName()); navCopy.NodeType() == xpath.TextNode && content != "" {
			element.Children = append(element.Children, &tableaupb.Child{
				Options: &tableaupb.ChildOptions{},
				Type: "", //TODO
				Name: "content",			
			})
			continue
		}
		// commentNode, documentNode and other meaningless nodes should be filtered
		if navCopy.NodeType() != xpath.ElementNode {
			continue
		}
		tagName := navCopy.LocalName()
		if _, exist := tagMap[tagName]; exist {
			continue
		}
		tagMap[tagName] = true
		newChild := &tableaupb.Child{
			Options: &tableaupb.ChildOptions{
				Name: tagName,
			},
			Type: tagName,
			Name: strcase.ToSnake(tagName),
			Element: &tableaupb.Element{},
		}
		if childList := xmlquery.Find(nav.Current(), tagName); len(childList) > 1 {
			newChild.Card = "repeated"
			newChild.Name = newChild.Name + "_list"
		}
		element.Children = append(element.Children, newChild)
		gen.parseNode(&navCopy, newChild.Element)
	}
	return nil
}