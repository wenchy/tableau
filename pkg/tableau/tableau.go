package tableau

import (
	"github.com/Wenchy/tableau/internal/converter"
)

// Tableaux is an alias type of converter.Tableaux.
type Tableaux = converter.Tableaux

func Convert(protoPackageName, inputPath, outputPath string) {
	tableaux := converter.Tableaux{ProtoPackageName: protoPackageName, InputPath: inputPath, OutputPath: outputPath}
	tableaux.Convert()
}

func NewTableaux(opts *Options) *Tableaux {
	opts.init()
	tbx := converter.Tableaux{
		ProtoPackageName:          opts.ProtoPackageName,
		InputPath:                 opts.InputPath,
		OutputPath:                opts.OutputPath,
		OutputFilenameAsSnakeCase: opts.OutputFilenameAsSnakeCase,
		OutputFormat:              converter.Format(opts.OutputFormat),
		OutputPretty:              opts.OutputPretty,
		LocationName:              opts.LocationName,
	}

	return &tbx
}
