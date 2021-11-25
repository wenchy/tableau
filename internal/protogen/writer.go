package protogen

import (
	"bufio"
	"fmt"
	"io"
)

type Writer struct {
	bw *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{bw: bufio.NewWriter(w)}
}

func (w *Writer) P(format string, a ...interface{}) {
	format += "\n"
	w.bw.WriteString(fmt.Sprintf(format, a...))
}

func (w *Writer) Flush() error {
	return w.bw.Flush()
}
