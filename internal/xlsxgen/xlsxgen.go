package xlsxgen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/iancoleman/strcase"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type metasheet struct {
	worksheet string // worksheet name
	namerow   int32  // exact row number of name at worksheet
	noterow   int32  // exact row number of description at wooksheet
	datarow   int32  // start row number of data
	transpose bool   // interchange the rows and columns
}

type Generator struct {
	ProtoPackage string // protobuf package name.
	InputDir         string // input dir of workbooks.
	OutputDir        string // output dir of generated protoconf files.

	metasheet metasheet // meta info of worksheet
}

var specialMessageMap = map[string]int{
	"google.protobuf.Timestamp": 1,
	"google.protobuf.Duration":  1,
}

type Cell struct {
	Name string
	Data string
}
type Row []Cell

func (gen *Generator) Generate() {
	err := os.RemoveAll(gen.OutputDir)
	if err != nil {
		panic(err)
	}
	// create output dir
	err = os.MkdirAll(gen.OutputDir, 0700)
	if err != nil {
		panic(err)
	}

	protoPackage := protoreflect.FullName(gen.ProtoPackage)
	protoregistry.GlobalFiles.RangeFilesByPackage(protoPackage, func(fd protoreflect.FileDescriptor) bool {
		atom.Log.Debugf("filepath: %s\n", fd.Path())
		opts := fd.Options().(*descriptorpb.FileOptions)
		workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(*tableaupb.WorkbookOptions)
		if workbook == nil {
			return true
		}

		atom.Log.Debugf("proto: %s => workbook: %s\n", fd.Path(), workbook)
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			md := msgs.Get(i)
			// atom.Log.Debugf("%s\n", md.FullName())
			opts := md.Options().(*descriptorpb.MessageOptions)
			worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
			if worksheet == nil {
				continue
			}
			atom.Log.Infof("generate: %s, message: %s#%s, worksheet: %s#%s", md.Name(), fd.Path(), md.Name(), workbook.Name, worksheet.Name)
			newMsg := dynamicpb.NewMessage(md)
			gen.export(newMsg)
		}
		return true
	})
}

// export the protomsg message.
func (gen *Generator) export(protomsg proto.Message) {
	md := protomsg.ProtoReflect().Descriptor()
	_, workbook := TestParseFileOptions(md.ParentFile())
	fmt.Println("==================", workbook)
	msgName, worksheet, namerow, noterow, datarow, transpose := TestParseMessageOptions(md)
	gen.metasheet.worksheet = worksheet
	gen.metasheet.namerow = namerow
	gen.metasheet.noterow = noterow
	gen.metasheet.datarow = datarow
	gen.metasheet.transpose = transpose

	row := make(Row, 0)
	gen.TestParseFieldOptions(md, &row, 0, "")
	fmt.Println("==================", msgName)

	filename := gen.OutputDir + workbook.Name
	var wb *excelize.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		wb = excelize.NewFile()
		t := time.Now()
		datetime := t.Format(time.RFC3339)
		err := wb.SetDocProps(&excelize.DocProperties{
			Category:       "category",
			ContentStatus:  "Draft",
			Created:        datetime,
			Creator:        "Tableau",
			Description:    "This file was created by Tableau",
			Identifier:     "xlsx",
			Keywords:       "Spreadsheet",
			LastModifiedBy: "Tableau",
			Modified:       datetime,
			Revision:       "0",
			Subject:        "Configuration",
			Title:          workbook.Name,
			Language:       "en-US",
			Version:        "1.0.0",
		})
		if err != nil {
			panic(err)
		}
		// The newly created workbook will by default contain a worksheet named `Sheet1`.
		wb.SetSheetName("Sheet1", worksheet)
		wb.SetDefaultFont("Courier")
	} else {
		fmt.Println("exist file: ", filename)
		wb, err = excelize.OpenFile(filename)
		if err != nil {
			panic(err)
		}
		wb.NewSheet(worksheet)
	}

	{
		style, err := wb.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:  "gradient",
				Color: []string{"#FFFFFF", "#E5E5E5"},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "top",
				WrapText:   true,
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Times New Roman",
				// Color:  "EEEEEEEE",
			},
			Border: []excelize.Border{
				{
					Type:  "top",
					Color: "EEEEEEEE",
					Style: 1,
				},
				{
					Type:  "bottom",
					Color: "EEEEEEEE",
					Style: 2,
				},
				{
					Type:  "left",
					Color: "EEEEEEEE",
					Style: 1,
				},
				{
					Type:  "right",
					Color: "EEEEEEEE",
					Style: 1,
				},
			},
		})
		if err != nil {
			panic(err)
		}
		wb.SetRowHeight(worksheet, 1, 50)
		for i, cell := range row {
			hanWidth := 1 * float64(getHanCount(cell.Name))
			letterWidth := 1 * float64(getLetterCount(cell.Name))
			digitWidth := 1 * float64(getDigitCount(cell.Name))
			width := hanWidth + letterWidth + digitWidth + 4.0
			// width := 2 * float64(utf8.RuneCountInString(cell.Name))
			colname, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				panic(err)
			}
			wb.SetColWidth(worksheet, colname, colname, width)

			axis, err := excelize.CoordinatesToCellName(i+1, 1)
			if err != nil {
				panic(err)
			}
			err = wb.SetCellValue(worksheet, axis, cell.Name)
			if err != nil {
				panic(err)
			}

			err = wb.AddComment(worksheet, axis, `{"author":"Tableau: ","text":"\n`+cell.Name+`, \nthis is a comment."}`)
			if err != nil {
				panic(err)
			}
			// set style
			wb.SetCellStyle(worksheet, axis, axis, style)
			if err != nil {
				panic(err)
			}
			atom.Log.Debugf("%s(%v) ", cell.Name, width)

			// test for validation
			// - min
			// - max
			// - droplist
			dataStartAxis, err := excelize.CoordinatesToCellName(i+1, 2)
			if err != nil {
				panic(err)
			}
			dataEndAxis, err := excelize.CoordinatesToCellName(i+1, 1000)
			if err != nil {
				panic(err)
			}

			if i == 0 {
				dataAxis, err := excelize.CoordinatesToCellName(i+1, 2)
				if err != nil {
					panic(err)
				}

				// unique key validation
				dv := excelize.NewDataValidation(true)
				dv.Sqref = dataStartAxis + ":" + dataEndAxis
				dv.Type = "custom"
				// dv.SetInput("Key", "Must be unique in this column")
				// NOTE(wenchyzhu): Five XML escape characters
				// "   &quot;
				// '   &apos;
				// <   &lt;
				// >   &gt;
				// &   &amp;
				//
				// `<formula1>=COUNTIF($A$2:$A$1000,A2)<2</formula1`
				//					||
				//					\/
				// `<formula1>=COUNTIF($A$2:$A$1000,A2)&lt;2</formula1`
				formula := fmt.Sprintf("=COUNTIF($A$2:$A$10000,%s)<2", dataAxis)
				dv.Formula1 = fmt.Sprintf("<formula1>%s</formula1>", escapeXml(formula))

				dv.SetError(excelize.DataValidationErrorStyleStop, "Error", "Key must be unique!")
				err = wb.AddDataValidation(worksheet, dv)
				if err != nil {
					panic(err)
				}
			} else if i == 1 {
				dv := excelize.NewDataValidation(true)
				dv.Sqref = dataStartAxis + ":" + dataEndAxis
				dv.SetDropList([]string{"1", "2", "3"})
				dv.SetInput("Options", "1: coin\n2: gem\n3: coupon")
				err := wb.AddDataValidation(worksheet, dv)
				if err != nil {
					panic(err)
				}
			} else if i == 2 {
				dv := excelize.NewDataValidation(true)
				dv.Sqref = dataStartAxis + ":" + dataEndAxis
				dv.SetRange(10, 20, excelize.DataValidationTypeWhole, excelize.DataValidationOperatorBetween)
				dv.SetError(excelize.DataValidationErrorStyleStop, "error title", "error body")
				err := wb.AddDataValidation(worksheet, dv)
				if err != nil {
					panic(err)
				}
			}

		}
		fmt.Println()
	}

	err := wb.SaveAs(filename)
	if err != nil {
		panic(err)
	}
}
func escapeXml(in string) string {
	var b bytes.Buffer
	err := xml.EscapeText(&b, []byte(in))
	if err != nil {
		panic(err)
	}
	return b.String()
}

func getHanCount(s string) int {
	count := 0
	for _, r := range []rune(s) {
		if unicode.Is(unicode.Han, r) {
			count++
		}
	}
	return count
}

func getLetterCount(s string) int {
	count := 0
	for _, r := range []rune(s) {
		if unicode.IsLetter(r) {
			count++
		}
	}
	return count
}

func getDigitCount(s string) int {
	count := 0
	for _, r := range []rune(s) {
		if unicode.IsDigit(r) {
			count++
		}
	}
	return count
}

// TestParseFileOptions is aimed to parse the options of a protobuf definition file.
func TestParseFileOptions(fd protoreflect.FileDescriptor) (string, *tableaupb.WorkbookOptions) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	protofile := string(fd.FullName())
	workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(*tableaupb.WorkbookOptions)
	atom.Log.Debugf("file:%s.proto, workbook:%s\n", protofile, workbook)
	return protofile, workbook
}

// TestParseMessageOptions is aimed to parse the options of a protobuf message.
func TestParseMessageOptions(md protoreflect.MessageDescriptor) (string, string, int32, int32, int32, bool) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgName := string(md.Name())
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)

	worksheetName := worksheet.Name
	namerow := worksheet.Namerow
	if worksheet.Namerow != 0 {
		namerow = 1 // default
	}
	noterow := worksheet.Noterow
	if noterow == 0 {
		noterow = 1 // default
	}
	datarow := worksheet.Datarow
	if datarow == 0 {
		datarow = 2 // default
	}
	transpose := worksheet.Transpose
	atom.Log.Debugf("message:%s, worksheetName:%s, namerow:%d, noterow:%d, datarow:%d, transpose:%v\n", msgName, worksheetName, namerow, noterow, datarow, transpose)
	return msgName, worksheetName, namerow, noterow, datarow, transpose
}

func getTabStr(depth int) string {
	tab := ""
	for i := 0; i < depth; i++ {
		tab += "\t"
	}
	return tab
}

// TestParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func (gen *Generator) TestParseFieldOptions(md protoreflect.MessageDescriptor, row *Row, depth int, prefix string) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
	worksheetName := ""
	if worksheet != nil {
		worksheetName = worksheet.Name
	}
	pkg := md.ParentFile().Package()
	atom.Log.Debugf("%s// %s, '%s', %v, %v, %v\n", getTabStr(depth), md.FullName(), worksheetName, md.IsMapEntry(), prefix, pkg)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if string(pkg) != gen.ProtoPackage && pkg != "google.protobuf" {
			atom.Log.Debugf("%s// no need to proces package: %v\n", getTabStr(depth), pkg)
			return
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
		}

		// default value
		name := strcase.ToCamel(string(fd.FullName().Name()))
		etype := tableaupb.Type_TYPE_DEFAULT
		key := ""
		layout := tableaupb.Layout_LAYOUT_DEFAULT
		sep := ""
		subsep := ""

		opts := fd.Options().(*descriptorpb.FieldOptions)
		field := proto.GetExtension(opts, tableaupb.E_Field).(*tableaupb.FieldOptions)
		if field != nil {
			name = field.Name
			etype = field.Type
			key = field.Key
			layout = field.Layout
			sep = field.Sep
			subsep = field.Subsep
		} else {
			// default processing
			if fd.IsList() {
				// truncate suffix `List` (CamelCase) corresponding to `_list` (snake_case)
				name = strings.TrimSuffix(name, "List")
			} else if fd.IsMap() {
				// truncate suffix `Map` (CamelCase) corresponding to `_map` (snake_case)
				// name = strings.TrimSuffix(name, "Map")
				name = ""
				key = "Key"
			}
		}
		if sep == "" {
			sep = ","
		}
		if subsep == "" {
			subsep = ":"
		}
		atom.Log.Debugf("%s%s(%v) %s(%s) %s = %d [(name) = \"%s\", (type) = %s, (key) = \"%s\", (layout) = \"%s\", (sep) = \"%s\"];",
			getTabStr(depth), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), prefix+name, etype.String(), key, layout.String(), sep)
		atom.Log.Debugw("field metadata",
			"tabs", depth,
			"cardinality", fd.Cardinality().String(),
			"isMap", fd.IsMap(),
			"kind", fd.Kind().String(),
			"msgName", msgName,
			"fullName", fd.FullName(),
			"number", fd.Number(),
			"name", prefix+name,
			"type", etype.String(),
			"key", key,
			"layout", layout.String(),
			"sep", sep,
		)
		if fd.IsMap() {
			valueFd := fd.MapValue()
			if etype == tableaupb.Type_TYPE_INCELL_MAP {
				if valueFd.Kind() == protoreflect.MessageKind {
					panic("in-cell map do not support value as message type")
				}
				fmt.Println("cell(FIELD_TYPE_CELL_MAP): ", prefix+name)
				*row = append(*row, Cell{Name: prefix + name})
			} else {
				if valueFd.Kind() == protoreflect.MessageKind {
					if layout == tableaupb.Layout_LAYOUT_HORIZONTAL {
						size := 2
						for i := 1; i <= size; i++ {
							// fmt.Println("cell: ", prefix+name+strconv.Itoa(i)+key)
							gen.TestParseFieldOptions(valueFd.Message(), row, depth+1, prefix+name+strconv.Itoa(i))
						}
					} else {
						// fmt.Println("cell: ", prefix+name+strconv.Itoa(i)+key)
						gen.TestParseFieldOptions(valueFd.Message(), row, depth+1, prefix+name)
					}
				} else {
					// value is scalar type
					key := "Key"     // deafult key name
					value := "Value" // deafult value name
					fmt.Println("cell(scalar map key): ", prefix+name+key)
					fmt.Println("cell(scalar map value): ", prefix+name+value)

					*row = append(*row, Cell{Name: prefix + name + key})
					*row = append(*row, Cell{Name: prefix + name + value})
				}
			}
		} else if fd.IsList() {
			if fd.Kind() == protoreflect.MessageKind {
				if layout == tableaupb.Layout_LAYOUT_VERTICAL {
					gen.TestParseFieldOptions(fd.Message(), row, depth+1, prefix+name)
				} else {
					size := 2
					for i := 1; i <= size; i++ {
						gen.TestParseFieldOptions(fd.Message(), row, depth+1, prefix+name+strconv.Itoa(i))
					}
				}
			} else {
				if etype == tableaupb.Type_TYPE_INCELL_LIST {
					fmt.Println("cell(FIELD_TYPE_CELL_LIST): ", prefix+name)
					*row = append(*row, Cell{Name: prefix + name})
				} else {
					panic(fmt.Sprintf("unknown list type: %v\n", etype))
				}
			}
		} else {
			if fd.Kind() == protoreflect.MessageKind {
				if etype == tableaupb.Type_TYPE_INCELL_STRUCT {
					fmt.Println("cell(FIELD_TYPE_CELL_MESSAGE): ", prefix+name)
					*row = append(*row, Cell{Name: prefix + name})
				} else {
					subMsgName := string(fd.Message().FullName())
					_, found := specialMessageMap[subMsgName]
					if found {
						fmt.Println("cell(special message): ", prefix+name)
						*row = append(*row, Cell{Name: prefix + name})
					} else {
						pkgName := fd.Message().ParentFile().Package()
						if string(pkgName) != gen.ProtoPackage {
							panic(fmt.Sprintf("unknown message %v in package %v", subMsgName, pkgName))
						}
						gen.TestParseFieldOptions(fd.Message(), row, depth+1, prefix+name)
					}
				}
			} else {
				fmt.Println("cell: ", prefix+name)
				*row = append(*row, Cell{Name: prefix + name})
			}
		}
	}
}
