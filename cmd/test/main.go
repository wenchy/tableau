package main

import (
	_ "github.com/Wenchy/tableau/cmd/test/testpb"
	"github.com/Wenchy/tableau/pkg/tableau"
)

func main() {
	tableau.Convert("test", "./testdata/")
}

/*
func parseActivity() {
	// fmt.Println("Hello, world.")
	// item := testpb.Item{}
	// Redact(item.ProtoReflect().Interface())

	conf := testpb.ActivityConf{
		ActivityMap: map[int32]*testpb.ActivityConf_Activity{
			1: &testpb.ActivityConf_Activity{
				ChapterMap: map[int32]*testpb.ActivityConf_Chapter{
					2: &testpb.ActivityConf_Chapter{
						SectionMap: map[int32]*testpb.ActivityConf_Row{
							3: &testpb.ActivityConf_Row{
								ActivityId: 1,
								ChapterId:  2,
								// ChapterDesc: "aha",
								SectionId: 3,
								// SectionDesc: "aha",
								Items: []*testpb.Item{
									&testpb.Item{
										Id:  1,
										Num: 2,
									},
									&testpb.Item{
										Id:  2,
										Num: 2,
									},
								},
								BeginTime: &timestamp.Timestamp{
									Seconds: 1596985188,
									Nanos:   1,
								},
								EndTime: &timestamp.Timestamp{
									Seconds: 1699985188,
									Nanos:   1,
								},
								Duration: &duration.Duration{
									Seconds: 3600,
									Nanos:   2,
								},
							},
						},
					},
				},
			},
		},
	}

	output, err := protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	var out bytes.Buffer
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()

	desc := conf.ActivityMap[1].ChapterMap[2].SectionMap[3].ActivityId
	fmt.Printf("desc: %v\n", desc)

	md := conf.ProtoReflect().Descriptor()

	ParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	ParseMessageOptions(md)
	fmt.Println("==================")
	ParseFieldOptions(md, 0)
	fmt.Println("==================")
}
// ParseFileOptions is aimed to parse the options of a protobuf definition file.
func ParseFileOptions(fd protoreflect.FileDescriptor) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(string)
	fmt.Printf("file:%s.proto, workbook:%s\n", fd.FullName(), workbook)
}

// ParseMessageOptions is aimed to parse the options of a protobuf message.
func ParseMessageOptions(md protoreflect.MessageDescriptor) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	captrow := proto.GetExtension(opts, tableaupb.E_Captrow).(int32)
	descrow := proto.GetExtension(opts, tableaupb.E_Descrow).(int32)
	datarow := proto.GetExtension(opts, tableaupb.E_Datarow).(int32)
	fmt.Printf("message:%s, worksheet:%s, captrow:%d, descrow:%d, datarow:%d\n", md.FullName(), worksheet, captrow, descrow, datarow)
}

func getTabStr(depth int) string {
	tab := ""
	for i := 0; i < depth; i++ {
		tab += "\t"
	}
	return tab
}

// ParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func ParseFieldOptions(md protoreflect.MessageDescriptor, depth int) {
	fmt.Printf("%s// %s\n", getTabStr(depth), md.FullName())
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.ParentFile().Package() != converter.TableauPackageName {
			return
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
			// fmt.Println(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			// ParseFieldOptions(fd.Message(), depth+1)
		}

		// if fd.IsList() {
		// 	fmt.Println("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }
		opts := fd.Options().(*descriptorpb.FieldOptions)
		caption := proto.GetExtension(opts, tableaupb.E_Caption).(string)
		etype := proto.GetExtension(opts, tableaupb.E_Type).(tableaupb.FieldType)
		key := proto.GetExtension(opts, tableaupb.E_Key).(string)
		fmt.Printf("%s%s %s(%s) %s = %d [(caption) = \"%s\", (type) = %s, (key) = \"%s\"];\n", getTabStr(depth), fd.Cardinality().String(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), caption, etype.String(), key)
		// fmt.Println(fd.ContainingMessage().FullName())

		if fd.Kind() == protoreflect.MessageKind {
			ParseFieldOptions(fd.Message(), depth+1)
		}
	}

	// m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
	// 	opts := fd.Options().(*descriptorpb.FieldOptions)
	// 	caption := proto.GetExtension(opts, tableaupb.E_Col).(string)
	// 	if caption != "" {
	// 		fmt.Println(fd.FullName().Name(), caption)
	// 		// fmt.Println(fd.ContainingMessage().FullName())
	// 	}
	// 	return true
	// })
}

func parseItem() {
	conf := testpb.ItemConf{
		ItemMap: map[int32]*testpb.ItemConf_Row{
			1: &testpb.ItemConf_Row{
				Id:     1,
				Name:   "金币",
				Desc:   "通用货币",
				IconId: 5001,
				Attributes: []*testpb.ItemConf_Attribute{
					&testpb.ItemConf_Attribute{
						Id:    1,
						Value: 2,
					},
					&testpb.ItemConf_Attribute{
						Id:    2,
						Value: 2,
					},
				},
				Effects: []int32{1, 2, 3},
				ExpiryTime: &timestamp.Timestamp{
					Seconds: 1596985188,
					Nanos:   1,
				},
			},
		},
	}

	output, err := protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	var out bytes.Buffer
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()

	// desc := conf.Activities[1].Chapters[2].Sections[3].Desc
	// fmt.Printf("desc: %s\n", desc)

	md := conf.ProtoReflect().Descriptor()
	msg := conf.ProtoReflect()
	_, workbook := converter.TestParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	_, worksheet, _, _, _ := converter.TestParseMessageOptions(md)
	fmt.Println("==================")
	sheet := converter.ReadSheet(converter.WorkbookRootDir+workbook, worksheet)
	// row 0: captrow
	// row 1 - MaxRow: datarow
	for nrow := 0; nrow < sheet.MaxRow; nrow++ {
		if nrow >= 1 {
			// row, err := sheet.Row(nrow)
			// if err != nil {
			// 	panic(err)
			// }
			kv := make(map[string]string)
			for i := 0; i < sheet.MaxCol; i++ {
				metaCell, err := sheet.Cell(0, i)
				if err != nil {
					panic(err)
				}
				key := metaCell.Value
				dataCell, err := sheet.Cell(nrow, i)
				if err != nil {
					panic(err)
				}
				value := dataCell.Value
				kv[key] = value
			}
			converter.TestParseFieldOptions(msg, kv, 0, "")
		}
		fmt.Println()
	}
	fmt.Println("==================")

	output, err = protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()
}
*/
