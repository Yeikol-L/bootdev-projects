package headers

import (
	"errors"
	"regexp"
	"strings"
)

const crlf = "\r\n"

var allowedCharacters = regexp.MustCompile("^[A-Za-z0-9!#$%&'*+.^_`|~-]+$")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	header, _, found := strings.Cut(string(data), crlf)
	if !found {
		return 0, false, nil
	}
	if len(header) == 0 {
		return 2, true, nil
	}
	fieldName, value, found := strings.Cut(header, ":")

	if strings.ContainsAny(fieldName, " ") || !allowedCharacters.MatchString(fieldName) {
		return 0, false, errors.New("Header malformed")
	}
	loweredFieldName := strings.ToLower(fieldName)
	if v, ok := h[loweredFieldName]; ok {
		h[loweredFieldName] = v + ", " + strings.Trim(value, " ")
	} else {
		h[loweredFieldName] = strings.Trim(value, " ")
	}
	return len(header) + len(crlf), false, nil
}
