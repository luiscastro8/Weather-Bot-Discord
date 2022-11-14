package mylogger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintln(t *testing.T) {
	t.Run("It should print [Info] at the start", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Println("hello world")
		assert.Contains(t, infoBuffer.String(), "[Info]")
	})

	t.Run("It should not print to the error writer", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Println("hello world")
		assert.Equal(t, "", errBuffer.String())
	})

	t.Run("It should print hello world", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Println("hello world")
		assert.Contains(t, infoBuffer.String(), "hello world")
	})

	t.Run("It should print hello world 2", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Println("hello", "world", 2)
		assert.Contains(t, infoBuffer.String(), "hello world 2")
	})
}

func TestErrorln(t *testing.T) {
	t.Run("It should print [Error] at the start", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Errorln("hello world")
		assert.Contains(t, errBuffer.String(), "[Error]")
	})

	t.Run("It should not print to the info writer", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Errorln("hello world")
		assert.Equal(t, "", infoBuffer.String())
	})

	t.Run("It should print hello world", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Errorln("hello world")
		assert.Contains(t, errBuffer.String(), "hello world")
	})

	t.Run("It should print hello world 2", func(t *testing.T) {
		var infoBuffer bytes.Buffer
		var errBuffer bytes.Buffer
		Init(&infoBuffer, &errBuffer)
		Errorln("hello", "world", 2)
		assert.Contains(t, errBuffer.String(), "hello world 2")
	})
}
