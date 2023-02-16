package compression

import (
	"fmt"
	"io"
)

type Extractor interface {
	Read(src []byte) (int, error)
	Close() error
}

func NewExtractor(input io.Reader, encoding string) (Extractor, error) {
	if isGzipEncoded(encoding) {
		return newExtractorGZIP(input)
	}

	return nil, fmt.Errorf("compression - NewExtractor: %w (%s)", ErrEncodingNotSupported, encoding)
}
