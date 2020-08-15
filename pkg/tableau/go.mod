module github.com/Wenchy/tableau/pkg/tableau

go 1.14

require (
	github.com/Wenchy/tableau/internal/converter v0.0.0
	github.com/Wenchy/tableau/tableaupb v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/Wenchy/tableau/internal/converter => ../../internal/converter

replace github.com/Wenchy/tableau/tableaupb => ../../pkg/tableaupb
