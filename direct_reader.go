package director

import "io"

type directReader struct {
	r io.Reader
	buf, rest []byte
}

func (d *directReader) Read(b []byte) (n int, err error) {
	if d.rest != nil {
		n = copy(b, d.rest)
		if n < len(d.rest) {
			d.rest = d.rest[n:]
			return
		} else {
			d.rest = nil
		}
	}
	r, err := d.r.Read(d.buf)
	if err != nil {
		return n + r, err
	}
	w := copy(b[n:], d.buf[:r])
	if (r - w) > 0 {
		d.rest = d.buf[w:r]
	}
	return w + n, nil
}

func NewDirectReader(r io.Reader, buf []byte) io.Reader {
	return &directReader{
		r: r,
		buf: buf,
	}
}

/*
type directReader struct {
	rd io.Reader
	buf []byte
	r, w int
}

func (r *directReader) Read(b []byte) (n int, err error) {
	if r.r > 0 {
		n = copy(b, r.buf[r.r:r.w])
		if n < (r.w - r.r) {
			r.r += n
			return
		} else {
			r.r = 0
		}
	}
	rx, err := r.rd.Read(r.buf)
	if err != nil {
		return n + rx, err
	}
	nx := copy(b[n:], r.buf[:rx])
	if (rx - nx) > 0 {
		r.r, r.w = nx, rx
	}
	return nx + n, nil
}

func NewDirectReader(r io.Reader, buf []byte) io.Reader {
	return &directReader{
		r: r,
		buf: buf,
	}
}
*/
