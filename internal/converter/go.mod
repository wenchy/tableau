module github.com/Wenchy/tableau/internal/converter

go 1.14

require (
	github.com/Wenchy/tableau/pkg/tableaupb v0.0.0-00010101000000-000000000000
	github.com/tealeg/xlsx/v3 v3.2.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/Wenchy/tableau/pkg/tableaupb => ../../pkg/tableaupb
