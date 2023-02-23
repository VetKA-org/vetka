package compression

import (
	"compress/gzip"
	"io"
	"strings"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

func isGzipEncoded(encoding string) bool {
	return strings.Contains(encoding, "gzip")
}

func newGzipDecoder(input io.Reader) (reader *gzip.Reader, err error) {
	reader, err = gzip.NewReader(input)

	return
}

type gzipEncoder struct {
	gin.ResponseWriter

	log *logger.Logger

	// Only Gzip is supported
	encoder *gzip.Writer

	// Supported values of Content-Type header
	supportedContent map[string]struct{}
}

func newGzipEncoder(log *logger.Logger, writer gin.ResponseWriter) *gzipEncoder {
	supportedContent := make(map[string]struct{})
	supportedContent["application/json; charset=utf-8"] = struct{}{}
	supportedContent["text/html; charset=utf-8"] = struct{}{}

	return &gzipEncoder{
		ResponseWriter:   writer,
		log:              log,
		supportedContent: supportedContent,
	}
}

func (g *gzipEncoder) Write(resp []byte) (int, error) {
	contentType := g.Header().Get("Content-Type")
	if _, ok := g.supportedContent[contentType]; !ok {
		g.log.Debug().Msg("Compression: not supported for " + contentType)

		return g.ResponseWriter.Write(resp)
	}

	if g.encoder == nil {
		encoder, err := gzip.NewWriterLevel(g.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			g.log.Error().Err(err).Msg("Compression - Write - gzip.NewWriterLevel")

			return 0, err
		}

		g.encoder = encoder
	}

	g.Header().Set("Content-Encoding", "gzip")

	return g.encoder.Write(resp)
}

func (g *gzipEncoder) Close() {
	if g.encoder != nil {
		g.encoder.Close()
	}
}
