package fastio

import (
	"io"
	"strconv"
)

const defaultWriterBufSize = 64 * 1024 // 64KB

type FastWriter struct {
	w   io.Writer
	buf []byte
	pos int
	err error

	autoFlush bool
	limit     int

	scratch []byte
}

type writerError struct {
	err error
}

func (we writerError) Error() string {
	return we.err.Error()
}

func (we writerError) Unwrap() error {
	return we.err
}

func (we writerError) Is(target error) bool {
	if target == nil {
		return we.err == nil
	}
	return we.err.Error() == target.Error()
}

func NewWriter(w io.Writer) *FastWriter {
	return &FastWriter{
		w:         w,
		buf:       make([]byte, defaultWriterBufSize),
		autoFlush: false,
		limit:     defaultWriterBufSize / 2,
		scratch:   make([]byte, 0, 64),
	}
}

func NewWriterWithAutoFlush(w io.Writer, limit int) *FastWriter {
	if limit <= 0 || limit > defaultWriterBufSize {
		limit = defaultWriterBufSize / 2
	}
	return &FastWriter{
		w:         w,
		buf:       make([]byte, defaultWriterBufSize),
		autoFlush: true,
		limit:     limit,
		scratch:   make([]byte, 0, 64),
	}

}

func (fw *FastWriter) Err() error {
	return fw.err
}

func (fw *FastWriter) Flush() error {
	if fw.err != nil {
		return fw.err
	}
	if fw.pos == 0 {
		return nil
	}
	n, err := fw.w.Write(fw.buf[:fw.pos])
	if err != nil {
		fw.err = writerError{err: err}
		return err
	}
	if n < fw.pos {
		fw.err = writerError{err: io.ErrShortWrite}
		return fw.err
	}
	fw.pos = 0
	return nil
}

func (fw *FastWriter) ensureSpace(n int) error {
	if fw.err != nil {
		return fw.err
	}
	if fw.pos == 0 {
		if err := fw.checkWriter(); err != nil {
			return err
		}
	}
	if n > len(fw.buf) {
		if err := fw.Flush(); err != nil {
			return err
		}
		_, err := fw.w.Write(fw.buf[:0])
		if err != nil {
			fw.err = writerError{err: err}
		}
		return err
	}
	if fw.pos+n > len(fw.buf) {
		if err := fw.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (fw *FastWriter) checkWriter() error {
	if fw.err != nil {
		return fw.err
	}
	if _, err := fw.w.Write(nil); err != nil {
		fw.err = writerError{err: err}
		return fw.err
	}
	return nil
}

func (fw *FastWriter) Write(p []byte) (int, error) {
	if fw.err != nil {
		return 0, fw.err
	}
	total := 0
	for len(p) > 0 {
		if err := fw.ensureSpace(len(p)); err != nil {
			return total, err
		}
		n := copy(fw.buf[fw.pos:], p)
		fw.pos += n
		p = p[n:]
		total += n
		if fw.autoFlush && fw.pos >= fw.limit {
			if err := fw.Flush(); err != nil {
				return total, err
			}
		}
	}
	return total, nil
}

func (fw *FastWriter) WriteByte(b byte) error {
	if err := fw.ensureSpace(1); err != nil {
		return err
	}
	fw.buf[fw.pos] = b
	fw.pos++
	if fw.autoFlush && fw.pos >= fw.limit {
		return fw.Flush()
	}
	return nil
}

func (fw *FastWriter) WriteBytes(b []byte) error {
	_, err := fw.Write(b)
	return err
}

func (fw *FastWriter) WriteString(s string) error {
	return fw.WriteBytes([]byte(s))
}

func (fw *FastWriter) WriteLine(s string) error {
	if err := fw.WriteString(s); err != nil {
		return err
	}
	return fw.WriteByte('\n')
}

func (fw *FastWriter) WriteInt(v int) error {
	fw.scratch = strconv.AppendInt(fw.scratch[:0], int64(v), 10)
	return fw.WriteBytes(fw.scratch)
}

func (fw *FastWriter) WriteUint64(v uint64) error {
	fw.scratch = strconv.AppendUint(fw.scratch[:0], v, 10)
	return fw.WriteBytes(fw.scratch)
}

func (fw *FastWriter) WriteFloat64(v float64, prec int) error {
	fw.scratch = strconv.AppendFloat(fw.scratch[:0], v, 'f', prec, 64)
	return fw.WriteBytes(fw.scratch)
}

func (fw *FastWriter) WriteInt64(v int64) error {
	fw.scratch = strconv.AppendInt(fw.scratch[:0], v, 10)
	return fw.WriteBytes(fw.scratch)
}
