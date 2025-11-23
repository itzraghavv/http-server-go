package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewHeaders() Headers {
	return make(Headers)
}

func TestParseHeader(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("   Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 26, n)
	assert.False(t, done)

	// valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Invalid Characters
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// uppercase
	headers = NewHeaders()
	data = []byte("CoNtEnT-TyPe: text/html\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Multiple values
	headers = NewHeaders()
	headers["set-person"] = "lane-loves-go"

	data = []byte("Set-Person: prime-loves-zig\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", headers["set-person"])
	assert.NoError(t, err)
	assert.Equal(t, 29, n)
	assert.False(t, done)
}
