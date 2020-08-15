package tableau

import (
	"github.com/Wenchy/tableau/internal/converter"
)

func Convert(ProtoPackageName, WorkbookRootDir string) {
	tableaux := converter.Tableaux{ProtoPackageName: ProtoPackageName, WorkbookRootDir: WorkbookRootDir}
	tableaux.Convert()
}
