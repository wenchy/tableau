package tableau

import (
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
	atom.Log.Infof("options inited: %+v", opts)
	if err := g.Generate(); err != nil {
		atom.Log.Errorf("generate failed: %+v", err)
		atom.Log.Panic(err)
	}
}

// Xlsx2Proto converts xlsx files (with meta header) to protoconf files.
func Xlsx2Proto(protoPackage, goPackage, indir, outdir string, setters ...options.Option) {
	opts := options.ParseOptions(setters...)
	g := protogen.Generator{
		ProtoPackage: protoPackage,
		GoPackage:    goPackage,
		InputDir:     indir,
		OutputDir:    outdir,
		Header:       opts.Header,
		Imports:      opts.Imports,
	}

	if err := g.Generate(); err != nil {
		atom.Log.Panic(err)
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
