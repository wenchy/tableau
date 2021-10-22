package main

import (
	"testing"

	"github.com/Wenchy/tableau"
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/options"
	_ "github.com/Wenchy/tableau/test/testpb"
)

func init() {
	atom.InitZap("debug")
}

func Test_ConvertExcelToJSON(t *testing.T) {
	tableau.NewTableaux(
		"test",
		"./testdata/",
		"./_output/json/",
		options.LogLevel("debug"),
	).Convert()

	// tableau.Generate("test", "./testdata/", "./_output/xlsx/")
}
