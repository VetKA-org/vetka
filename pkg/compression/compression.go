package compression

import (
	"fmt"
	"io"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Decoder interface {
	Read(src []byte) (int, error)

	Close() error
}

func NewDecoder(input io.Reader, encoding string) (Decoder, error) {
	if isGzipEncoded(encoding) {
		return newGzipDecoder(input)
	}

	return nil, fmt.Errorf("compression - NewDecoder: %w (%s)", ErrEncodingNotSupported, encoding)
}

type Encoder interface {
	gin.ResponseWriter

	Write(resp []byte) (int, error)

	Close()
}

func NewEncoder(log *logger.Logger, writer gin.ResponseWriter, encoding string) (Encoder, error) {
	if isGzipEncoded(encoding) {
		return newGzipEncoder(log, writer), nil
	}

	return nil, fmt.Errorf("compression - NewEncoder: %w (%s)", ErrEncodingNotSupported, encoding)
}
