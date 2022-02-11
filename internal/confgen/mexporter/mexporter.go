// mexporter is the message exporter package, which can export one
// single message to different formts: JSON, Prototext, and Protowire.
package mexporter

import (
	"io/ioutil"
	"path/filepath"

	"github.com/Wenchy/tableau/options"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type messageExporter struct {
	name      string
	msg       proto.Message
	outputOpt *options.OutputOption
	outputDir string
}

func New(name string, msg proto.Message, outputDir string, outputOpt *options.OutputOption) *messageExporter {
	return &messageExporter{
		name:      name,
		msg:       msg,
		outputDir: outputDir,
		outputOpt: outputOpt,
	}
}

func (x *messageExporter) Export() error {
	filename := x.name
	if x.outputOpt.FilenameAsSnakeCase {
		filename = strcase.ToSnake(x.name)
	}

	var out []byte
	var err error
	switch x.outputOpt.Format {
	case options.JSON:
		filename += ".json"
		out, err = x.marshalToJSON()
		if err != nil {
			return errors.Wrapf(err, "failed to export %s to JSON", x.name)
		}

	case options.Prototext:
		filename += ".prototext"
		out, err = x.marshalToPrototext()
		if err != nil {
			return errors.Wrapf(err, "failed to export %s to prototext", x.name)
		}
	case options.Protowire:
		filename += ".protowire"
		out, err = x.marshalToProtowire()
		if err != nil {
			return errors.Wrapf(err, "failed to export %s to protowire", x.name)
		}
	default:
		return errors.Errorf("unknown output format: %v", x.outputOpt.Format)
	}

	fpath := filepath.Join(x.outputDir, filename)
	err = ioutil.WriteFile(fpath, out, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write file: %s", fpath)
	}
	// out.WriteTo(os.Stdout)
	return nil
}

func (x *messageExporter) marshalToJSON() (out []byte, err error) {
	if x.outputOpt.Pretty {
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "    ",
			EmitUnpopulated: x.outputOpt.EmitUnpopulated,
		}
		return opts.Marshal(x.msg)
	}
	return protojson.Marshal(x.msg)
}

func (x *messageExporter) marshalToPrototext() (out []byte, err error) {
	if x.outputOpt.Pretty {
		opts := prototext.MarshalOptions{
			Multiline: true,
			Indent:    "    ",
		}
		return opts.Marshal(x.msg)
	}
	return prototext.Marshal(x.msg)
}

func (x *messageExporter) marshalToProtowire() (out []byte, err error) {
	return proto.Marshal(x.msg)
}
