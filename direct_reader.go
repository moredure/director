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

func (d *directReader) Read(b []byte) (n int, err error) {
	if d.r != d.w {
		n = copy(b, d.buf[d.w:d.r])
		if n < d.r - d.w {
			d.w += n
			return
		}
	}
	d.r, err = d.rd.Read(d.buf)
	if err != nil {
		return n + d.r, err
	}
	d.w = copy(b[n:], d.buf[:d.r])
	return d.w + n, nil
}

func NewDirectReader(r io.Reader, buf []byte) io.Reader {
	return &directReader{
		r: r,
		buf: buf,
	}
}
*/
