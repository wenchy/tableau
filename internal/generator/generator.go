package generator

import (
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/Wenchy/tableau/pkg/tableaupb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// type Format int

// // file format
// const (
// 	Proto Format = 0
// 	Xlsx         = 1
// )

type metasheet struct {
	worksheet string // worksheet name
	captrow   int32  // exact row number of caption at worksheet
	descrow   int32  // exact row number of description at wooksheet
	datarow   int32  // start row number of data
	transpose bool   // interchange the rows and columns
}

type Generator struct {
	ProtoPackageName string // protobuf package name.
	InputPath        string // root dir of workbooks.
	OutputPath       string // output path of generated files.
	// OutputFormat     Format // output format: Proto, xlsx.

	metasheet metasheet // meta info of worksheet
}

var specialMessageMap = map[string]int{
	"google.protobuf.Timestamp": 1,
	"google.protobuf.Duration":  1,
}

type Cell struct {
	Caption string
	Data    string
}
type Row []Cell

func (gen *Generator) Generate() {
	err := os.RemoveAll(gen.OutputPath)
	if err != nil {
		panic(err)
	}
	// create oupput dir
	err = os.MkdirAll(gen.OutputPath, 0700)
	if err != nil {
		panic(err)
	}

	protoPackage := protoreflect.FullName(gen.ProtoPackageName)
	protoregistry.GlobalFiles.RangeFilesByPackage(protoPackage, func(fd protoreflect.FileDescriptor) bool {
		fmt.Printf("filepath: %s\n", fd.Path())
		opts := fd.Options().(*descriptorpb.FileOptions)
		workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(string)
		if workbook == "" {
			return true
		}

		fmt.Printf("proto: %s => workbook: %s\n", fd.Path(), workbook)
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			md := msgs.Get(i)
			// fmt.Printf("%s\n", md.FullName())
			opts := md.Options().(*descriptorpb.MessageOptions)
			worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
			if worksheet != "" {
				fmt.Printf("message: %s, worksheet: %s\n", md.FullName(), worksheet)
			}
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
	msgName, worksheet, captrow, descrow, datarow, transpose := TestParseMessageOptions(md)
	gen.metasheet.worksheet = worksheet
	gen.metasheet.captrow = captrow
	gen.metasheet.descrow = descrow
	gen.metasheet.datarow = datarow
	gen.metasheet.transpose = transpose

	row := make(Row, 0)
	gen.TestParseFieldOptions(md, &row, 0, "")
	fmt.Println("==================", msgName)

	filename := gen.OutputPath + workbook
	var wb *excelize.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		wb = excelize.NewFile()
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

		for i, cell := range row {
			hanWidth := 1 * float64(getHanCount(cell.Caption))
			letterWidth := 1 * float64(getLetterCount(cell.Caption))
			digitWidth := 1 * float64(getDigitCount(cell.Caption))
			width := hanWidth + letterWidth + digitWidth + 4.0
			// width := 2 * float64(utf8.RuneCountInString(cell.Caption))
			colname, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				panic(err)
			}
			wb.SetColWidth(worksheet, colname, colname, width)
			wb.SetRowHeight(worksheet, i+1, 50)

			axis, err := excelize.CoordinatesToCellName(i+1, 1)
			if err != nil {
				panic(err)
			}
			err = wb.SetCellValue(worksheet, axis, cell.Caption)
			if err != nil {
				panic(err)
			}

			// set style
			wb.SetCellStyle(worksheet, axis, axis, style)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s(%v) ", cell.Caption, width)

		}
		fmt.Println()
	}

	err := wb.SaveAs(filename)
	if err != nil {
		panic(err)
	}
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
func TestParseFileOptions(fd protoreflect.FileDescriptor) (string, string) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	protofile := string(fd.FullName())
	workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(string)
	fmt.Printf("file:%s.proto, workbook:%s\n", protofile, workbook)
	return protofile, workbook
}

// TestParseMessageOptions is aimed to parse the options of a protobuf message.
func TestParseMessageOptions(md protoreflect.MessageDescriptor) (string, string, int32, int32, int32, bool) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgName := string(md.Name())
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	captrow := proto.GetExtension(opts, tableaupb.E_Captrow).(int32)
	if captrow == 0 {
		captrow = 1 // default
	}
	descrow := proto.GetExtension(opts, tableaupb.E_Descrow).(int32)
	if descrow == 0 {
		descrow = 1 // default
	}
	datarow := proto.GetExtension(opts, tableaupb.E_Datarow).(int32)
	if datarow == 0 {
		datarow = 2 // default
	}
	transpose := proto.GetExtension(opts, tableaupb.E_Transpose).(bool)
	fmt.Printf("message:%s, worksheet:%s, captrow:%d, descrow:%d, datarow:%d, transpose:%v\n", msgName, worksheet, captrow, descrow, datarow, transpose)
	return msgName, worksheet, captrow, descrow, datarow, transpose
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
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	pkg := md.ParentFile().Package()
	fmt.Printf("%s// %s, '%s', %v, %v, %v\n", getTabStr(depth), md.FullName(), worksheet, md.IsMapEntry(), prefix, pkg)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if string(pkg) != gen.ProtoPackageName && pkg != "google.protobuf" {
			fmt.Printf("%s// no need to proces package: %v\n", getTabStr(depth), pkg)
			return
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
		}

		opts := fd.Options().(*descriptorpb.FieldOptions)
		caption := proto.GetExtension(opts, tableaupb.E_Caption).(string)
		etype := proto.GetExtension(opts, tableaupb.E_Type).(tableaupb.FieldType)
		key := proto.GetExtension(opts, tableaupb.E_Key).(string)
		layout := proto.GetExtension(opts, tableaupb.E_Layout).(tableaupb.CompositeLayout)
		sep := proto.GetExtension(opts, tableaupb.E_Sep).(string)
		if sep == "" {
			sep = ","
		}
		subsep := proto.GetExtension(opts, tableaupb.E_Subsep).(string)
		if subsep == "" {
			subsep = ":"
		}
		fmt.Printf("%s%s(%v) %s(%s) %s = %d [(caption) = \"%s\", (type) = %s, (key) = \"%s\", (layout) = \"%s\", (sep) = \"%s\"];\n",
			getTabStr(depth), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), prefix+caption, etype.String(), layout.String(), key, sep)
		if fd.IsMap() {
			valueFd := fd.MapValue()
			if etype == tableaupb.FieldType_FIELD_TYPE_CELL_MAP {
				if valueFd.Kind() == protoreflect.MessageKind {
					panic("in-cell map do not support value as message type")
				}
				fmt.Println("cell(FIELD_TYPE_CELL_MAP): ", prefix+caption)
				*row = append(*row, Cell{Caption: prefix + caption})
			} else {
				if valueFd.Kind() == protoreflect.MessageKind {
					if layout == tableaupb.CompositeLayout_COMPOSITE_LAYOUT_HORIZONTAL {
						size := 2
						for i := 1; i <= size; i++ {
							// fmt.Println("cell: ", prefix+caption+strconv.Itoa(i)+key)
							gen.TestParseFieldOptions(valueFd.Message(), row, depth+1, prefix+caption+strconv.Itoa(i))
						}
					} else {
						// fmt.Println("cell: ", prefix+caption+strconv.Itoa(i)+key)
						gen.TestParseFieldOptions(valueFd.Message(), row, depth+1, prefix+caption)
					}
				} else {
					// value is scalar type
					key := "Key"     // deafult key caption
					value := "Value" // deafult value caption
					fmt.Println("cell(scalar map key): ", prefix+caption+key)
					fmt.Println("cell(scalar map value): ", prefix+caption+value)

					*row = append(*row, Cell{Caption: prefix + caption + key})
					*row = append(*row, Cell{Caption: prefix + caption + value})
				}
			}
		} else if fd.IsList() {
			if fd.Kind() == protoreflect.MessageKind {
				if layout == tableaupb.CompositeLayout_COMPOSITE_LAYOUT_VERTICAL {
					gen.TestParseFieldOptions(fd.Message(), row, depth+1, prefix+caption)
				} else {
					size := 2
					for i := 1; i <= size; i++ {
						gen.TestParseFieldOptions(fd.Message(), row, depth+1, prefix+caption+strconv.Itoa(i))
					}
				}
			} else {
				if etype == tableaupb.FieldType_FIELD_TYPE_CELL_LIST {
					fmt.Println("cell(FIELD_TYPE_CELL_LIST): ", prefix+caption)
					*row = append(*row, Cell{Caption: prefix + caption})
				} else {
					panic(fmt.Sprintf("unknown list type: %v\n", etype))
				}
			}
		} else {
			if fd.Kind() == protoreflect.MessageKind {
				if etype == tableaupb.FieldType_FIELD_TYPE_CELL_MESSAGE {
					fmt.Println("cell(FIELD_TYPE_CELL_MESSAGE): ", prefix+caption)
					*row = append(*row, Cell{Caption: prefix + caption})
				} else {
					subMsgName := string(fd.Message().FullName())
					_, found := specialMessageMap[subMsgName]
					if found {
						fmt.Println("cell(special message): ", prefix+caption)
						*row = append(*row, Cell{Caption: prefix + caption})
					} else {
						pkgName := fd.Message().ParentFile().Package()
						if string(pkgName) != gen.ProtoPackageName {
							panic(fmt.Sprintf("unknown message %v in package %v", subMsgName, pkgName))
						}
						gen.TestParseFieldOptions(fd.Message(), row, depth+1, prefix+caption)
					}
				}
			} else {
				fmt.Println("cell: ", prefix+caption)
				*row = append(*row, Cell{Caption: prefix + caption})
			}
		}
	}
}
