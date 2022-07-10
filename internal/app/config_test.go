package app

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadJson(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	var err error
	var r io.Reader

	r = bytes.NewReader([]byte("not json"))
	err = cfg.LoadJson(r)
	assert.Error(err)
	assert.ErrorContains(err, "invalid character")

	r = bytes.NewReader([]byte("{}"))
	err = cfg.LoadJson(r)
	assert.NoError(err)
}

func TestConfig_SaveJson(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	var err error

	w := new(strings.Builder)
	err = cfg.SaveJson(w)
	assert.NoError(err)
	assert.True(len(w.String()) > 0)
}
