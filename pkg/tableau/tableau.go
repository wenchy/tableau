package tableau

import (
	"github.com/Wenchy/tableau/internal/converter"
)

func Convert(protoPackageName, inputPath, outputPath string, filenameAsSnakeCase bool) {
	tableaux := converter.Tableaux{ProtoPackageName: protoPackageName, InputPath: inputPath, OutputPath: outputPath, FilenameAsSnakeCase: filenameAsSnakeCase}
	tableaux.Convert()
}
