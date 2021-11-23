package main

import (
	"testing"

	"github.com/Wenchy/tableau"
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/options"
)

func init() {
	atom.InitZap("debug")
}
func Test_GenerateProtoconf(t *testing.T) {
	tableau.Xlsx2Protoconf(
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
			},
		),
	)
}
