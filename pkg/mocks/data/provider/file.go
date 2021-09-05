package provider

import (
	"io"
)

type File struct {
	name string
	buf  []byte
	ptr  int64
}


func (m *File) Close() error {
	return nil
}

func (m *File) Read(p []byte) (n int, err error) {
	n, err = m.ReadAt(p, m.ptr)
	m.ptr += int64(n)
	return n, err
}

func (m *File) Seek(offset int64, whence int) (int64, error) {
	return m.ptr, nil
}

func (m *File) ReadAt(p []byte, off int64) (n int, err error) {
	if n = copy(p, m.buf[off:]); n == 0 {
		return n, io.EOF
	} else {
		return n, nil
	}
}
