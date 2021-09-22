package director

import (
	"errors"
	"io"
)

// adaptation of bufio.Reader for directio
type DirectReader struct {
	rd   io.Reader
	buf  []byte
	r, w int
	err  error
}

var errNegativeRead = errors.New("err negative read")

func (d *DirectReader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		if d.Buffered() > 0 {
			return 0, nil
		}
		return 0, d.readErr()
	}
	if d.r == d.w {
		if d.err != nil {
			return 0, d.readErr()
		}
		d.r = 0
		d.w = 0
		n, d.err = d.rd.Read(d.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, d.readErr()
		}
		d.w += n
	}

	n = copy(p, d.buf[d.r:d.w])
	d.r += n
	return n, nil
}

func (d *DirectReader) Buffered() int { return d.w - d.r }

func (d *DirectReader) readErr() error {
	err := d.err
	d.err = nil
	return err
}

func NewDirectReader(r io.Reader, buf []byte) *DirectReader {
	return &DirectReader{
		rd:  r,
		buf: buf,
	}
}
