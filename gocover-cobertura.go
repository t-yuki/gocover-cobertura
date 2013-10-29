package main

import (
	"encoding/xml"
	"fmt"
	"go/build"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.Writer) {
	profiles, err := ParseProfiles(in)
	if err != nil {
		panic("Can't parse profiles")
	}

	srcDirs := build.Default.SrcDirs()
	sources := make([]Source, len(srcDirs))
	for i, dir := range srcDirs {
		sources[i] = Source{dir}
	}

	coverage := Coverage{Sources: sources, Packages: nil, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)}
	convertProfilesToPackages(profiles, &coverage)

	fmt.Fprintf(out, xml.Header)
	fmt.Fprintf(out, "<!DOCTYPE coverage SYSTEM \"http://cobertura.sourceforge.net/xml/coverage-03.dtd\">\n")

	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	err = encoder.Encode(coverage)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(out)
}

func convertProfilesToPackages(profiles []*Profile, coverage *Coverage) {
	coverage.Packages = []Package{}
}

func stripKnownSources(sources []Source, fileName string) string {
	for _, source := range sources {
		prefix := source.Path
		prefix = strings.TrimSuffix(prefix, string(os.PathSeparator)) + string(os.PathSeparator)
		if strings.HasPrefix(fileName, prefix) {
			return strings.TrimPrefix(fileName, prefix)
		}
	}
	return fileName
}
