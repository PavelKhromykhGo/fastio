// Package fastio предоставляет быстрые функции ввода/вывода,
// основанные на буферизации и низкоуровневом парсинге.
// FastReader и FastWriter ускоряют работу со stdin/stdout и файлами
// по сравнению со стандартными fmt.Fscan / fmt.Fprint.
package fastio

import (
	"errors"
	"io"
	"strconv"
)

const defaultReaderBufSize = 64 * 1024 // 64KB

// FastReader — быстрый буферизованный ридер.
//
// Он обеспечивает:
//   - минимальное количество аллокаций;
//   - методы для чтения примитивов: NextInt, NextInt64, NextUint64,
//     NextFloat64, NextWord, NextLine;
//   - совместимость с любым io.Reader (stdin, файл, сокет);
//   - ручное управление ошибками через Err().
//
// FastReader не является потокобезопасным.
type FastReader struct {
	r   io.Reader
	buf []byte
	pos int
	n   int
	err error
}

// NewReader создает FastReader поверх существующего io.Reader.
// Буфер создаётся размером 64 KB.
//
// Используется для быстрого чтения из stdin, файла или сетевого потока.
func NewReader(r io.Reader) *FastReader {
	return &FastReader{
		r:   r,
		buf: make([]byte, defaultReaderBufSize),
	}
}

// Err возвращает первую возникшую ошибку (включая io.EOF).
// Если ошибка уже произошла, дальнейшее чтение недоступно.
func (fr *FastReader) Err() error {
	return fr.err
}

func (fr *FastReader) fill() {
	if fr.err != nil {
		return
	}
	n, err := fr.r.Read(fr.buf)
	if n < 0 {
		n = 0
	}
	fr.n = n
	fr.pos = 0
	if err != nil {
		fr.err = err
	}
}

// ReadByte читает один байт из внутреннего буфера.
// При необходимости буфер автоматически заполняется.
//
// В случае достижения конца файла возвращает io.EOF.
func (fr *FastReader) ReadByte() (byte, error) {
	if err := fr.ensureData(); err != nil {
		return 0, err
	}

	b := fr.buf[fr.pos]
	fr.pos++

	if fr.pos >= fr.n {
		// Try to read ahead to determine the correct terminal error state.
		if fr.err == nil {
			fr.fill()
		}
		if fr.n == 0 {
			if fr.err == nil {
				fr.err = io.EOF
			}
			return b, fr.err
		}
	}

	if fr.pos >= fr.n && fr.err != nil {
		return b, fr.err
	}

	return b, nil
}

// PeekByte возвращает следующий байт, не продвигая позицию в буфере.
// Полезно для lookahead-парсинга.
//
// В случае EOF возвращается ошибка io.EOF.
func (fr *FastReader) PeekByte() (byte, error) {
	if err := fr.ensureData(); err != nil {
		return 0, err
	}

	return fr.buf[fr.pos], nil
}

func (fr *FastReader) ensureData() error {
	if fr.err != nil && !(fr.err == io.EOF && fr.pos < fr.n) {
		return fr.err
	}

	if fr.pos >= fr.n {
		fr.fill()
		if fr.n == 0 {
			if fr.err == nil {
				fr.err = io.EOF
			}
			return fr.err
		}
	}

	return nil
}

// SkipSpaces пропускает пробельные символы: пробелы, \n, \r, \t.
// Используется перед парсингом чисел и слов.
//
// Если пробелы находятся в конце файла — возвращает io.EOF.
func (fr *FastReader) SkipSpaces() error {
	for {
		b, err := fr.PeekByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			_, _ = fr.ReadByte()
			continue
		}
		return nil
	}
}

// NextWord читает последовательность непробельных символов.
// Используется для токенизации входа.
//
// В случае отсутствия данных возвращает io.EOF.
func (fr *FastReader) NextWord() (string, error) {
	if err := fr.SkipSpaces(); err != nil {
		if errors.Is(err, io.EOF) {
			return "", io.EOF
		}
		return "", err
	}

	var buf []byte
	for {
		b, err := fr.PeekByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(buf) == 0 {
					return "", io.EOF
				}
				return string(buf), nil
			}
			return "", err
		}
		if b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			break
		}
		_, _ = fr.ReadByte()
		buf = append(buf, b)
	}
	if len(buf) == 0 {
		return "", io.EOF
	}
	return string(buf), nil
}

// NextInt читает целое число типа int (со знаком).
// Формат поддерживает ведущие пробелы, знак '+' или '-'.
//
// В случае отсутствия цифр возвращает ошибку.
func (fr *FastReader) NextInt() (int, error) {
	if err := fr.SkipSpaces(); err != nil {
		return 0, err
	}

	sign := 1
	b, err := fr.PeekByte()
	if err != nil {
		return 0, err
	}
	if b == '-' {
		sign = -1
		_, _ = fr.ReadByte()
	} else if b == '+' {
		_, _ = fr.ReadByte()
	}

	var val int
	digitsRead := 0

	for {
		b, err = fr.PeekByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		if b < '0' || b > '9' {
			break
		}
		_, _ = fr.ReadByte()
		val = val*10 + int(b-'0')
		digitsRead++
	}

	if digitsRead == 0 {
		return 0, errors.New("fastio: NextInt: no digits found")
	}
	return sign * val, nil
}

// NextInt64 читает 64-битное целое число со знаком.
// Работает аналогично NextInt, но возвращает int64.
func (fr *FastReader) NextInt64() (int64, error) {
	if err := fr.SkipSpaces(); err != nil {
		return 0, err
	}

	sign := int64(1)
	b, err := fr.PeekByte()
	if err != nil {
		return 0, err
	}
	if b == '-' {
		sign = -1
		_, _ = fr.ReadByte()
	} else if b == '+' {
		_, _ = fr.ReadByte()
	}

	var val int64
	digitsRead := 0

	for {
		b, err = fr.PeekByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		if b < '0' || b > '9' {
			break
		}
		_, _ = fr.ReadByte()
		val = val*10 + int64(b-'0')
		digitsRead++
	}
	if digitsRead == 0 {
		return 0, errors.New("fastio: NextInt64: no digits found")
	}
	return sign * val, nil
}

// NextUint64 читает беззнаковое целое число.
// Допускается ведущий '+' перед числом.
//
// В случае отсутствия цифр возвращает ошибку.
func (fr *FastReader) NextUint64() (uint64, error) {
	if err := fr.SkipSpaces(); err != nil {
		return 0, err
	}

	b, err := fr.PeekByte()
	if err != nil {
		return 0, err
	}

	if b == '+' {
		_, _ = fr.ReadByte()
	}

	var val uint64
	digitsRead := 0

	for {
		b, err = fr.PeekByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		if b < '0' || b > '9' {
			break
		}
		_, _ = fr.ReadByte()
		val = val*10 + uint64(b-'0')
		digitsRead++
	}
	if digitsRead == 0 {
		return 0, errors.New("fastio: NextUint64: no digits found")
	}
	return val, nil
}

// NextFloat64 читает число в формате float64.
// Поддерживает: целые, дробные, экспоненциальные ("1e9") форматы.
// Реализовано через NextWord + strconv.ParseFloat.
func (fr *FastReader) NextFloat64() (float64, error) {
	token, err := fr.NextWord()
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// NextLine читает строку до символа '\n'.
// Символ переноса строки не включается в результат.
// CRLF ("\r\n") приводится к обычному LF.
//
// В случае пустого оставшегося ввода возвращает io.EOF.
func (fr *FastReader) NextLine() (string, error) {
	var buf []byte

	for {
		b, err := fr.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(buf) == 0 {
					return "", io.EOF
				}
				break
			}
			return "", err
		}
		if b == '\n' {
			break
		}
		buf = append(buf, b)
	}
	if len(buf) > 0 && buf[len(buf)-1] == '\r' {
		buf = buf[:len(buf)-1]
	}

	return string(buf), nil
}
