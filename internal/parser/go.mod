module github.com/claesp/verto/internal/parser

go 1.22.5

require github.com/claesp/verto/internal/importer v0.0.0

require github.com/claesp/verto/internal/types v0.0.0 // indirect

replace (
	github.com/claesp/verto/internal/importer => ../importer/
	github.com/claesp/verto/internal/types => ../types/
)
