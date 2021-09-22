package director

import (
	"errors"
	"io"
)

// adaptation of bufio.Reader for directio
type directReader struct {
	rd io.Reader
	buf []byte
	r, w int
	err error
}

var errNegativeRead = errors.New("err negative read")

func (b *directReader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		if b.Buffered() > 0 {
			return 0, nil
		}
		return 0, b.readErr()
	}
	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		b.r = 0
		b.w = 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}

func (b *directReader) Buffered() int { return b.w - b.r }

func (b *directReader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

func NewDirectReader(r io.Reader, buf []byte) io.Reader {
	return &directReader{
		rd: r,
		buf: buf,
	}
}
