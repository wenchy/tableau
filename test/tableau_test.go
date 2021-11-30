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

func Test_Xlsx2Proto(t *testing.T) {
	tableau.Xlsx2Proto(
		"test",
		"github.com/Wenchy/tableau/cmd/test/testpb",
		"./testdata",
		"./protoconf",
		options.Header(
			&options.HeaderOption{
				Namerow: 1,
				Typerow: 2,
				Noterow: 3,
				Datarow: 5,
			}),
		options.Imports(
			[]string{
				"common.proto",
				"time.proto",
			},
		),
	)
}

func Test_Xlsx2JSON(t *testing.T) {
	tableau.Xlsx2Conf(
		"test",
		"./testdata/",
		"./_output/json/",
		options.LogLevel("debug"),
	)
	// tableau.Generate("test", "./testdata/", "./_output/xlsx/")
}
