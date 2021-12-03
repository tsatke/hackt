package ui

import (
	"bytes"
	"io"
)

var _ io.Writer = (*EditorContent)(nil)
var _ io.WriterAt = (*EditorContent)(nil)

type EditorContent struct {
	data []byte

	linesCache [][]byte
}

func NewBuffer(data []byte) *EditorContent {
	return &EditorContent{
		data: data,
	}
}

func (c *EditorContent) Bytes() []byte {
	return c.data
}

func (c *EditorContent) Lines() [][]byte {
	if c.linesCache == nil {
		c.linesCache = bytes.Split(c.data, []byte{'\n'})
	}
	return c.linesCache
}

func (c *EditorContent) invalidateLinesCache() {
	c.linesCache = nil
}

func (c *EditorContent) Write(p []byte) (n int, err error) {
	c.invalidateLinesCache()

	c.data = append(c.data, p...)
	return len(p), nil
}

func (c *EditorContent) WriteAt(p []byte, off int64) (n int, err error) {
	c.invalidateLinesCache()

	c.ensureSize(off + int64(len(p)))
	copy(c.data[off:off+int64(len(p))], p)
	return len(p), nil
}

func (c *EditorContent) InsertAt(p []byte, off int64) {
	c.invalidateLinesCache()

	c.data = append(c.data[:off], append(p, c.data[off:]...)...)
}

func (c *EditorContent) DeleteAt(len int64, off int64) {
	c.invalidateLinesCache()

	c.data = append(c.data[:off], c.data[off+len:]...)
}

func (c *EditorContent) Copy() *EditorContent {
	newData := make([]byte, len(c.data))
	copy(newData[:len(c.data)], c.data)
	return NewBuffer(newData)
}

func (c *EditorContent) ReplaceAllBytes(toReplace []byte, replacement []byte) {
	c.data = bytes.ReplaceAll(c.data, toReplace, replacement)
}

func (c *EditorContent) ensureSize(size int64) {
	if int64(len(c.data)) < size {
		newSlice := make([]byte, size)
		copy(newSlice[:len(c.data)], c.data)
	}
}
