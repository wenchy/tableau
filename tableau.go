package tableau

import (
	"os"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/confgen"
	"github.com/Wenchy/tableau/internal/protogen"
	"github.com/Wenchy/tableau/internal/xlsxgen"
	"github.com/Wenchy/tableau/options"
)

// Xlsx2Conf converts xlsx files (with meta header) to different formatted configuration files.
// Supported formats: json, prototext, and protowire.
func Xlsx2Conf(protoPackage, indir, outdir string, setters ...options.Option) {
	opts := options.ParseOptions(setters...)
	g := confgen.Generator{
		ProtoPackage: protoPackage,
		LocationName: opts.LocationName,
		InputDir:     indir,
		OutputDir:    outdir,
		Output:       opts.Output,
	}
	atom.InitZap(opts.LogLevel)
	atom.Log.Debugf("options inited: %+v, header: %+v, output: %+v", opts, opts.Header, opts.Output)
	if err := g.Generate(opts.Workbook, opts.Worksheet); err != nil {
		atom.Log.Errorf("generate failed: %+v", err)
		os.Exit(-1)
	}
}

// Xlsx2Proto converts xlsx files (with meta header) to protoconf files.
func Xlsx2Proto(protoPackage, goPackage, indir, outdir string, setters ...options.Option) {
	opts := options.ParseOptions(setters...)
	g := protogen.Generator{
		ProtoPackage: protoPackage,
		GoPackage:    goPackage,
		LocationName: opts.LocationName,
		InputDir:     indir,
		OutputDir:    outdir,

		FilenameWithSubdirPrefix: opts.Output.FilenameWithSubdirPrefix,
		FilenameSuffix:           opts.Output.FilenameSuffix,

		Header:  opts.Header,
		Imports: opts.Imports,
	}
	atom.InitZap(opts.LogLevel)
	atom.Log.Debugf("options inited: %+v, header: %+v, output: %+v", opts, opts.Header, opts.Output)
	if err := g.Generate(); err != nil {
		atom.Log.Errorf("generate failed: %+v", err)
		os.Exit(-1)
	}
}

// Proto2Xlsx converts protoconf files to xlsx files (with meta header).
func Proto2Xlsx(protoPackage, indir, outdir string) {
	g := xlsxgen.Generator{
		ProtoPackage: protoPackage,
		InputDir:     indir,
		OutputDir:    outdir,
	}
	g.Generate()
}
