package bundles

import (
	"flag"

	"github.com/rs/zerolog"
)

// TODO: The idea is instead of doing $ref parsing and dealing with file vs. http, we use our artifact:// and hydrate at build time.
// tl;dr, we dont use $ref for our own references in our bundle.yaml files

// TODO: build files...
// [ ] metadata.json
// [ ] inputs, connections, artifacts JSON
func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	bundleFilePath := flag.String("input", "./bundle.yaml", "Path to bundle file")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// [metadata, artifactJson, inputJson, connectionJson] =
	ParseBundle(*bundleFilePath)
	// each to file
}
