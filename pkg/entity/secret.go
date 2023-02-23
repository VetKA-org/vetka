package entity

import (
	"regexp"
	"strings"
)

type Secret string

func (s Secret) String() string {
	return strings.Repeat("*", len(s))
}

type SecretURI string

var _URISecrets = regexp.MustCompile(`(://).*:.*(@)`)

func (u SecretURI) String() string {
	return string(_URISecrets.ReplaceAll([]byte(u), []byte("$1*****:*****$2")))
}
