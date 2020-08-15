package tableau

import (
	"github.com/Wenchy/tableau/internal/converter"
)

func Convert(ProtoPackageName, WorkbookRootDir, OutputPath string) {
	tableaux := converter.Tableaux{ProtoPackageName: ProtoPackageName, WorkbookRootDir: WorkbookRootDir, OutputPath: OutputPath}
	tableaux.Convert()
}
