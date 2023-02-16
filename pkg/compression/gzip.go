package compression

import (
	"compress/gzip"
	"io"
	"strings"
)

func isGzipEncoded(encoding string) bool {
	return strings.Contains(encoding, "gzip")
}

func newExtractorGZIP(input io.Reader) (reader *gzip.Reader, err error) {
	reader, err = gzip.NewReader(input)

	return
}
