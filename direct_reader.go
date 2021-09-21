package director

import "io"

type directReader struct {
	r io.Reader
	buf, rest []byte
}

func (r *directReader) Read(b []byte) (n int, err error) {
	if r.rest != nil {
		n = copy(b, r.rest)
		if n < len(r.rest) {
			r.rest = r.rest[n:]
			return
		} else {
			r.rest = nil
		}
	}
	rx, err := r.r.Read(r.buf)
	if err != nil {
		return n + rx, err
	}
	nx := copy(b[n:], r.buf[:rx])
	if (rx - nx) > 0 {
		r.rest = r.buf[nx:rx]
	}
	return nx + n, nil
}

func NewDirectReader(r io.Reader, buf []byte) io.Reader {
	return &directReader{
		r: r,
		buf: buf,
	}
}
