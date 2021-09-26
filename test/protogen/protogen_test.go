package main

import (
	"testing"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/protogen"
)

func init() {
	atom.InitZap("debug")
}
func Test_Generate(t *testing.T) {
	generator := protogen.Generator{
		ProtoPackage: "test",
		GoPackage:    "github.com/Wenchy/tableau/cmd/test/testpb",
		InputDir:     "./testdata",
		OutputDir:    "./protoconf",
	}
	generator.Generate()
}
