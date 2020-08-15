module github.com/Wenchy/tableau/cmd/test

go 1.14

require (
	github.com/Wenchy/tableau/pkg/tableau v0.0.0-00010101000000-000000000000
	github.com/Wenchy/tableau/cmd/test/testpb v0.0.0-00010101000000-000000000000
)

replace github.com/Wenchy/tableau/pkg/tableau => ../../pkg/tableau

replace github.com/Wenchy/tableau/pkg/tableaupb => ../../pkg/tableaupb

replace github.com/Wenchy/tableau/cmd/test/testpb => ./testpb

replace github.com/Wenchy/tableau/internal/converter => ../../internal/converter
