module github.com/Wenchy/tableau/pkg/tableau

go 1.14

require github.com/Wenchy/tableau/internal/converter v0.1.0

replace github.com/Wenchy/tableau/internal/converter => ../../internal/converter

replace github.com/Wenchy/tableau/pkg/tableaupb => ../../pkg/tableaupb
