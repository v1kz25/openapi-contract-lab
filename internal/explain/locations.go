package explain

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/v1kz25/openapi-contract-lab/internal/oasdiff"
)

func sourceLocation(side string, source *oasdiff.Source, input string) Location {
	location := Location{Side: side, File: source.File, Line: source.Line, Column: source.Column}
	if input == "" || !filepath.IsAbs(input) || source.File == "" || source.Line < 0 || source.Column < 0 {
		return location
	}
	// Link only a source that matches the known local input. In Windows the
	// pinned engine may omit the drive letter; never guess by basename alone.
	inputPath := filepath.Clean(input)
	match := matchesInput(source.File, inputPath)
	// The pinned engine also returns URI-escaped paths for names containing #,
	// spaces or non-ASCII characters. Decode once, then match the full input.
	if !match {
		if decoded, err := url.PathUnescape(source.File); err == nil {
			match = matchesInput(decoded, inputPath)
		}
	}
	if !match {
		return location
	}
	path := filepath.ToSlash(inputPath)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	link := url.URL{Scheme: "file", Path: path}
	if source.Line > 0 {
		link.Fragment = fmt.Sprintf("L%d", source.Line)
	}
	location.URL = link.String()
	return location
}

func matchesInput(source, input string) bool {
	path := filepath.Clean(source)
	if input == path {
		return true
	}
	if runtime.GOOS == "windows" && (strings.EqualFold(input, path) || strings.EqualFold(strings.TrimPrefix(input, filepath.VolumeName(input)), path)) {
		return true
	}
	if !filepath.IsAbs(source) {
		if absolute, err := filepath.Abs(source); err == nil {
			return input == absolute || (runtime.GOOS == "windows" && strings.EqualFold(input, absolute))
		}
	}
	return false
}
