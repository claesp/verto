module github.com/claesp/verto/verto-cli

go 1.22.5

require (
	github.com/claesp/verto/internal/importer v0.0.0
	github.com/claesp/verto/internal/parser v0.0.0
)

require github.com/claesp/verto/internal/types v0.0.0 // indirect

replace (
	github.com/claesp/verto/internal/importer => ../../internal/importer/
	github.com/claesp/verto/internal/parser => ../../internal/parser/
	github.com/claesp/verto/internal/types => ../../internal/types/
)
