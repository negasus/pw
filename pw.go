package pw

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var (
	ColorReset  = []byte("\033[0m")
	ColorRed    = []byte("\033[31m")
	ColorGreen  = []byte("\033[32m")
	ColorYellow = []byte("\033[33m")
	ColorBlue   = []byte("\033[34m")
	ColorPurple = []byte("\033[35m")
	ColorCyan   = []byte("\033[36m")
	ColorGray   = []byte("\033[37m")
	ColorWhite  = []byte("\033[97m")
)

func FPrint(w io.Writer, data []byte) error {
	j := New(data)
	j.w = w
	return j.Render()
}

func Print(data []byte) error {
	return New(data).Render()
}

type byteReader interface {
	io.Reader
	io.ByteReader
}

type PrettyWriter struct {
	ch  byte
	t   int
	r   byteReader
	w   io.Writer
	tab []byte
	err error

	pos int

	colorNumber []byte
	colorString []byte
	colorBool   []byte
	colorNull   []byte

	plainArray bool
}

func New(data []byte) *PrettyWriter {
	pw := &PrettyWriter{
		r:           bytes.NewReader(data),
		w:           os.Stdout,
		tab:         []byte("    "),
		colorNumber: ColorCyan,
		colorString: ColorGreen,
		colorBool:   ColorYellow,
		colorNull:   ColorBlue,
	}

	return pw
}

func (pw *PrettyWriter) Render() error {
	if !pw.next() {
		return pw.error()
	}
	if pw.ch == '{' || pw.ch == '[' {
		return pw.renderJSON()
	}
	if pw.ch == '<' {
		return pw.renderXML()
	}

	_, err := io.Copy(pw.w, pw.r)
	return err
}

func (pw *PrettyWriter) SetWriter(w io.Writer) {
	pw.w = w
}

func (pw *PrettyWriter) SetTab(data []byte) {
	pw.tab = data
}

func (pw *PrettyWriter) Reset(data []byte) {
	pw.r = bytes.NewReader(data)
	pw.pos = 0
	pw.t = 0
	pw.err = nil
	pw.ch = 0
}

func (pw *PrettyWriter) nextNoSpace() bool {
	for {
		pw.pos++
		b, err := pw.r.ReadByte()
		if err != nil {
			pw.err = err
			return false
		}
		pw.ch = b
		return true
	}
}

func (pw *PrettyWriter) next() bool {
	for {
		pw.pos++
		b, err := pw.r.ReadByte()
		if err != nil {
			pw.err = err
			return false
		}
		if isSpace(b) {
			continue
		}
		pw.ch = b
		return true
	}
}

func (pw *PrettyWriter) error() error {
	return fmt.Errorf("error %v at position %d", pw.err, pw.pos)
}

func (pw *PrettyWriter) withColor(color []byte, f func() bool) bool {
	pw.writeByte(color...)
	r := f()
	pw.writeByte(ColorReset...)
	return r
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
