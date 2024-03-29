package main

import (
	"testing"

	"github.com/Wenchy/tableau"
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/importer"
	"github.com/Wenchy/tableau/options"
	_ "github.com/Wenchy/tableau/test/testpb/excel"
	_ "github.com/Wenchy/tableau/test/testpb/xml"
)

func init() {
	atom.InitZap("debug")
}

func Test_Excel2Proto(t *testing.T) {
	tableau.Excel2Proto(
		"testexcel",
		"github.com/Wenchy/tableau/cmd/test/testpb/excel",
		"./testdata/excel",
		"./protoconf/excel",
		options.Header(
			&options.HeaderOption{
				Namerow: 1,
				Typerow: 2,
				Noterow: 3,
				Datarow: 5,

				Nameline: 2,
				Typeline: 2,
			}),
		options.Imports(
			[]string{
				"common.proto",
				"time.proto",
			},
		),
		options.Output(
			&options.OutputOption{
				FilenameSuffix:           "_conf",
				FilenameWithSubdirPrefix: true,
			},
		),
		options.LogLevel("debug"),
	)
}

func Test_Excel2JSON(t *testing.T) {
	tableau.Excel2Conf(
		"testexcel",
		"./testdata/excel",
		"./_output/json/",
		options.LogLevel("debug"),
	)
}

func Test_Excel2JSON_Select(t *testing.T) {
	tableau.Excel2Conf(
		"testexcel",
		"./testdata/excel",
		"./_output/json/",
		options.LogLevel("debug"),
		// options.Workbook("hero/Test.xlsx"),
		// options.Workbook("./hero/Test.xlsx"),
		options.Workbook(".\\hero\\Test.xlsx"),
		options.Worksheet("Hero"),
	)
}

func Test_Excel2CSV(t *testing.T) {
	paths := []string{
		"./testdata/excel/Test.xlsx",
		"./testdata/excel/hero/Test.xlsx",
	}
	for _, path := range paths {
		imp := importer.NewExcelImporter(path, nil, nil, true)
		err := imp.ExportCSV()
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_CSV2Excel(t *testing.T) {
	paths := []string{
		"./testdata/excel/Test#Activity.csv",
		"./testdata/excel/Test#Reward.csv",
		"./testdata/excel/Test#Exchange.csv",
		"./testdata/excel/Test#Match.csv",
		"./testdata/excel/Test#Loader.csv",
		"./testdata/excel/Test#@TABLEAU.csv",
		"./testdata/excel/Test#Sheet2.csv",

		"./testdata/excel/hero/Test#Hero.csv",
		"./testdata/excel/hero/Test#@TABLEAU.csv",
	}
	for _, path := range paths {
		imp := importer.NewCSVImporter(path)
		err := imp.ExportExcel()
		if err != nil {
			t.Errorf("%+v", err)
		}
	}
}

func Test_XML2Proto(t *testing.T) {
	tableau.XML2Proto(
		"testxml",
		"github.com/Wenchy/tableau/cmd/test/testpb/xml",
		"./testdata/xml",
		"./protoconf/xml",
		options.Imports(
			[]string{
				"cs_com_def.proto",
			},
		),
		options.Input(
			&options.InputOption{
				Format: options.XML,
			},
		),
	)
}

func Test_XML2JSON(t *testing.T) {
	tableau.XML2Conf(
		"testxml",
		"./testdata/xml",
		"./_output/json",
		options.LogLevel("debug"),
		options.Input(
			&options.InputOption{
				Format: options.XML,
			},
		),
	)
	// tableau.Generate("test", "./testdata/", "./_output/xml/")
}
