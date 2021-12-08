package confgen

import (
	"os"
	"regexp"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

var newlineRegex *regexp.Regexp

func init() {
	newlineRegex = regexp.MustCompile(`\r?\n?`)
}

func clearNewline(s string) string {
	return newlineRegex.ReplaceAllString(s, "")
}

type Generator struct {
	ProtoPackage string // protobuf package name.
	LocationName string // Location represents the collection of time offsets in use in a geographical area. Default is "Asia/Shanghai".
	InputDir     string // input dir of workbooks.
	OutputDir    string // output dir of generated files.

	Output *options.OutputOption // output settings.
}

var specialMessageMap = map[string]int{
	"google.protobuf.Timestamp": 1,
	"google.protobuf.Duration":  1,
}

func (gen *Generator) Generate() (err error) {
	// create output dir
	err = os.MkdirAll(gen.OutputDir, 0700)
	if err != nil {
		return errors.WithMessagef(err, "failed to create output dir: %s", gen.OutputDir)
	}

	protoregistry.GlobalFiles.RangeFilesByPackage(
		protoreflect.FullName(gen.ProtoPackage),
		func(fd protoreflect.FileDescriptor) bool {
			// atom.Log.Debugf("filepath: %s", fd.Path())
			opts := fd.Options().(*descriptorpb.FileOptions)
			workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(*tableaupb.WorkbookOptions)
			if workbook == nil {
				return true
			}
			// atom.Log.Debugf("proto: %s, workbook %s", fd.Path(), workbook)
			msgs := fd.Messages()
			for i := 0; i < msgs.Len(); i++ {
				md := msgs.Get(i)
				// atom.Log.Debugf("%s", md.FullName())
				opts := md.Options().(*descriptorpb.MessageOptions)
				worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
				if worksheet == nil {
					continue
				}
				atom.Log.Infof("generate: %s@%s <-> %s@%s", fd.Path(), md.Name(), workbook.Name, worksheet.Name)
				newMsg := dynamicpb.NewMessage(md)
				// gen.export(newMsg)
				x := &sheetParser{
					ProtoPackage: gen.ProtoPackage,
					LocationName: gen.LocationName,
					InputDir:     gen.InputDir,
					OutputDir:    gen.OutputDir,
					Output:       gen.Output,

					protomsg: newMsg,
				}
				err = x.Export()
				if err != nil {
					// Due to closure, this err will be returned by func Generate().
					return false
				}
			}
			return true
		})
	return err
}
