package protogen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"

	"google.golang.org/protobuf/proto"
	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/iancoleman/strcase"
	"github.com/emirpasic/gods/sets/treeset"
)

type XmlGenerator struct {
	ProtoPackage string // protobuf package name.
	GoPackage    string // golang package name.
	InputDir     string // input dir of workbooks.
	OutputDir    string // output dir of generated protoconf files.
	Imports      []string // imported common proto file paths

	Xml *options.XmlOption // xml generation settings
	writer *bufio.Writer
	childMap map[string]*tableaupb.Child
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
		gen.childMap = make(map[string]*tableaupb.Child)
		gen.parseNode(xmlquery.CreateXPathNavigator(n), xml.xml.Root, "")
		if err := gen.exportXml(xml.xml); err != nil {
			atom.Log.Panic(err)
		}
	}
}

func (gen *XmlGenerator) parseNode(nav *xmlquery.NodeNavigator, element *tableaupb.Element, prefix string) error {	
	element.Options = &tableaupb.ElementOptions{
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
			keyValue := keyNode.InnerText()
			fmt.Println(keyValue)
			element.Options.Key = attr.Value
			element.KeyType = "uint32" //TODO
		case "desc":
		default:
			attrName := attr.Name.Local
			// attrValue := attr.Value
			newAttr := &tableaupb.Attr{
				Options: &tableaupb.AttrOptions{
					Name: attrName,
					Default: "0", //TODO
				},
				Type: "int", //TODO
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
	// fmt.Println(element)

	// createMetaStruct create a meta struct by some pre-defined tag
	tagMap := make(map[string]bool)
	createMetaStruct := func(metaNode *xmlquery.Node) {
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
		gen.parseNode(metaNav, newChild.Element, fmt.Sprintf("%s/%s", prefix, metaNode.Parent.Data))
		gen.childMap[fmt.Sprintf("%s/%s/%s", prefix, metaNode.Parent.Data, metaNav.LocalName())] = newChild
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
		newChild := &tableaupb.Child{
			Options: &tableaupb.ChildOptions{
				Name: tagName,
			},
			Type: tagName,
			Name: strcase.ToSnake(tagName),
			Element: &tableaupb.Element{},
		}
		if keyTag, childList := strings.Split(element.Options.Key, ".")[0], xmlquery.Find(nav.Current(), tagName); len(childList) > 1 || tagName == keyTag {			
			newChild.Card = "repeated"
			newChild.Name = newChild.Name + "_list"
		}
		fatherPath := fmt.Sprintf("%s/%s", prefix, nav.Current().Data)
		gen.parseNode(&navCopy, newChild.Element, fatherPath)
		
		// overwrite previous meta struct if necessary
		curPath := fatherPath + "/" + tagName
		if child, exist := gen.childMap[curPath]; exist {
			if child.Card == "" && newChild.Card == "repeated" {
				child.Card = "repeated"
				child.Name = child.Name + "_list"
			}
			if newChild.Element != nil {
				if child.Element == nil {
					child.Element = newChild.Element					
				} else {
					attrMap := make(map[string]bool)
					for _, attr := range child.Element.Attrs {
						attrMap[attr.Name] = true
					}
					for _, attr := range newChild.Element.Attrs {
						if _, exist := attrMap[attr.Name]; !exist {
							child.Element.Attrs = append(child.Element.Attrs, attr)
						}
					}
					childMap := make(map[string]bool)
					for _, c := range child.Element.Children {
						childMap[c.Options.Name] = true
					}
					for _, c := range newChild.Element.Children {
						if _, exist := childMap[c.Options.Name]; !exist {
							child.Element.Children = append(child.Element.Children, c)
						}
					}
				}
			}
		} else {
			gen.childMap[curPath] = newChild
			element.Children = append(element.Children, newChild)
		}
	}
	return nil
}

func (gen *XmlGenerator) exportXml(xml *tableaupb.Xml) error {
	atom.Log.Debug(proto.Marshal(xml))
	path := filepath.Join(gen.OutputDir, xml.Name+".proto")
	atom.Log.Debugf("output: %s", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gen.writer = bufio.NewWriter(f)
	defer gen.writer.Flush()

	gen.writer.WriteString("syntax = \"proto3\";\n")
	gen.writer.WriteString(fmt.Sprintf("package %s;\n", gen.ProtoPackage))
	gen.writer.WriteString(fmt.Sprintf("option go_package = \"%s\";\n", gen.GoPackage))
	gen.writer.WriteString("\n")

	// keep the elements ordered by sheet name
	set := treeset.NewWithStringComparator()
	for key := range xml.Imports {
		set.Add(key)
	}
	for _, key := range set.Values() {
		gen.writer.WriteString(fmt.Sprintf("import \"%s\";\n", key))
	}
	gen.writer.WriteString("\n")
	gen.writer.WriteString(fmt.Sprintf("option (tableau.xml) = {%s};\n", genPrototext(xml.Options)))
	gen.writer.WriteString("\n")

	gen.exportElement(xml.Root, 0)

	return nil
}

func (gen *XmlGenerator) exportElement(element *tableaupb.Element, depth int) error {
	gen.writer.WriteString(indent(depth) + fmt.Sprintf("message %s {\n", element.Name))
	gen.writer.WriteString(indent(depth) + fmt.Sprintf("  option (tableau.element) = {%s};\n", genPrototext(element.Options)))
	gen.writer.WriteString("\n")
	tagid := 0
	// generate attributes
	for _, attr := range element.Attrs {
		tagid++
		attrLine := fmt.Sprintf("%s %s = %d [(tableau.attr) = {%s}];\n", attr.Type, attr.Name, tagid, genPrototext(attr.Options))
		if attr.Card == "repeated" {
			attrLine = "repeated " + attrLine
		}
		gen.writer.WriteString(indent(depth) + "  " + attrLine)
	}
	if len(element.Attrs) > 0 {
		gen.writer.WriteString("\n")
	}
	// generate child elements
	if element.Options.Key != "" {
		tagid++
		gen.writer.WriteString(indent(depth) + fmt.Sprintf("  map<%s, int32> %s_map = %d;\n", element.KeyType, strcase.ToSnake(strings.Split(element.Options.Key, ".")[0]), tagid))
	}
	for _, child := range element.Children {
		tagid++
		childLine := fmt.Sprintf("%s %s = %d [(tableau.child) = {%s}];\n", child.Type, child.Name, tagid, genPrototext(child.Options))
		if child.Card == "repeated" {
			childLine = "repeated " + childLine
		}
		gen.writer.WriteString(indent(depth) + "  " + childLine)
	}
	// generate child messages
	for _, child := range element.Children {
		if child.Element != nil {
			gen.writer.WriteString("\n")
			gen.exportElement(child.Element, depth + 1)
		}		
	}
	gen.writer.WriteString(indent(depth) + "}\n")
	
	return nil
}