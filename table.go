package pw

import (
	"bytes"
	"io"
)

type Table struct {
	header     []string
	data       [][]string
	withHeader bool

	cols int

	maxLen map[int]int
}

func NewTable() *Table {
	return &Table{
		maxLen: make(map[int]int),
	}
}

func (t *Table) SetHeader(values ...string) *Table {
	t.header = values
	t.withHeader = true

	for i, v := range values {
		if len(v) > t.maxLen[i] {
			t.maxLen[i] = len(v)
		}
	}

	if t.cols < len(values) {
		t.cols = len(values)
	}

	return t
}

func (t *Table) AddLine(values ...string) *Table {
	t.data = append(t.data, values)

	for i, v := range values {
		if len(v) > t.maxLen[i] {
			t.maxLen[i] = len(v)
		}
	}

	if t.cols < len(values) {
		t.cols = len(values)
	}

	return t
}

func (t *Table) Render(w io.Writer) {
	if t.withHeader {
		t.renderSeparator(w)
		t.renderLine(w, t.header)
	}
	t.renderSeparator(w)
	for _, line := range t.data {
		t.renderLine(w, line)
	}
	t.renderSeparator(w)

}

func (t *Table) renderSeparator(w io.Writer) {
	for i := 0; i < t.cols; i++ {
		w.Write([]byte("+"))
		w.Write(bytes.Repeat([]byte("-"), t.maxLen[i]+2))
	}
	w.Write([]byte("+\n"))
}

func (t *Table) renderLine(w io.Writer, line []string) {
	for i, v := range line {
		if i > 0 {
			w.Write([]byte(" "))
		}
		w.Write([]byte("| "))
		w.Write([]byte(v))
		w.Write(bytes.Repeat([]byte(" "), t.maxLen[i]-len(v)))
	}

	if len(line) < t.cols {
		for i := len(line); i < t.cols; i++ {
			w.Write([]byte(" |"))
			w.Write(bytes.Repeat([]byte(" "), t.maxLen[i]+1))
		}
	}

	w.Write([]byte(" |\n"))
}
