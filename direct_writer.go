package director

import (
	"io"
)

type DirectWriter struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}

func NewDirectWriter(w io.Writer, buf []byte) *DirectWriter {
	return &DirectWriter{
		buf: buf,
		wr:  w,
	}
}

func (b *DirectWriter) Size() int { return len(b.buf) }

func (b *DirectWriter) Reset(w io.Writer) {
	b.err = nil
	b.n = 0
	b.wr = w
}

func (b *DirectWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf)
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

func (b *DirectWriter) Available() int { return len(b.buf) - b.n }

func (b *DirectWriter) Buffered() int { return b.n }

func (b *DirectWriter) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.wr.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}