package director

import "io"

type directReader struct {
	rd io.Reader
	buf []byte
	r, w int
}

func (d *directReader) Read(b []byte) (n int, err error) {
	if d.r == d.w {
		d.r, err = d.rd.Read(d.buf)
		if err != nil {
			return 0, err
		}
		d.w = copy(b, d.buf[:d.r])
		return d.w, nil
	}
	n = copy(b, d.buf[d.w:d.r])
	d.w += n
	return n, nil
}

func NewDirectReader(r io.Reader, buf []byte) io.Reader {
	return &directReader{
		rd: r,
		buf: buf,
	}
}
