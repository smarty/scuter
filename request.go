//go:build goexperiment.jsonv2

package scuter

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ReadUint64Header parses the first header value corresponding with key as a uint64.
func ReadUint64Header(headers http.Header, key string) uint64 {
	if values, contains := headers[key]; contains && len(values) > 0 {
		parsed, _ := strconv.ParseUint(values[0], 10, 64)
		return parsed
	}
	return 0
}

// ReadTimeHeader parses the first header value corresponding with key as a time.Time, else returns the Epoch time.
func ReadTimeHeader(headers http.Header, format, key string) (zero time.Time) {
	if values, contains := headers[key]; contains && len(values) > 0 {
		parsed, _ := time.Parse(format, values[0])
		return parsed
	}
	return zero
}

// ReadPathElement returns the path element immediately following the supplied label.
func ReadPathElement(rawPath, label string) string {
	found := false
	for element := range strings.SplitSeq(rawPath, "/") {
		if found == true {
			return element
		}
		if element == label {
			found = true
		}
	}
	return ""
}

// ReadNumericPathElement returns the path element immediately following the supplied label, parsed as a uint64
func ReadNumericPathElement(rawPath, label string) uint64 {
	raw := ReadPathElement(rawPath, label)
	value, _ := strconv.ParseUint(raw, 10, 64)
	return value
}
