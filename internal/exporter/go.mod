module github.com/claesp/verto/internal/exporter

go 1.22.5

require (
	github.com/claesp/verto/internal/types v0.0.0
)

replace (
	github.com/claesp/verto/internal/types => ../types/
)
