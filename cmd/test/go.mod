module github.com/Wenchy/tableau/cmd/test

go 1.14

require (
	github.com/Wenchy/tableau/converter v0.0.0-20200811131850-045cfae77a14
	github.com/Wenchy/tableau/tableaupb v0.0.0-00010101000000-000000000000
	github.com/Wenchy/tableau/testpb v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.2
	github.com/tealeg/xlsx/v3 v3.2.0 // indirect
	google.golang.org/protobuf v1.25.0
)

replace github.com/Wenchy/tableau/converter => ../../internal/converter

replace github.com/Wenchy/tableau/testpb => ./testpb

replace github.com/Wenchy/tableau/tableaupb => ../../pkg/tableaupb
