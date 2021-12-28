package main

import (
	"github.com/Wenchy/tableau"
	_ "github.com/Wenchy/tableau/cmd/test/testpb"
	"github.com/Wenchy/tableau/options"
)

func main() {
	tableau.Xlsx2Conf(
		"test",
		"./testdata/",
		"./_output/json/",
		options.LogLevel("debug"),
	)

	// tableau.Proto2Xlsx("test", "./testdata/", "./_output/xlsx/")
}
