package main

import (
	"testing"

	"github.com/Wenchy/tableau"
	"github.com/Wenchy/tableau/internal/atom"
	_ "github.com/Wenchy/tableau/test/testpb"
)

func init() {
	atom.InitZap("debug")
}

func Test_ConvertExcelToJSON(t *testing.T) {

	// tableau.Convert("test", "./testdata/", "./output/")
	tbx := tableau.NewTableaux(&tableau.Options{
		ProtoPackage: "test",
		InputDir:     "./testdata/",
		OutputDir:    "./_output/json/",
		// OutputFilenameAsSnakeCase: false,
		OutputFormat:    tableau.JSON,
		OutputPretty:    true,
		EmitUnpopulated: true,
		LogLevel:        "debug",
	})
	tbx.Convert()

	// tableau.Generate("test", "./testdata/", "./_output/xlsx/")
}
