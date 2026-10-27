package runner

import (
	"bufio"
	"encoding/binary"
	"io"
	"unicode/utf16"
	"unicode/utf8"
)

// byteOrderMark is the invisible character some programs put at the start
// of UTF-16 text. It is dropped from the output.
const byteOrderMark rune = 0xFEFF

// readableOutput returns a program's output as UTF-8 text. Most programs
// write UTF-8 or plain ASCII, which passes through unchanged. Some Windows
// tools, such as SFC, write UTF-16 when their output goes to another program
// instead of a console window. UTF-16 gives itself away in the first two
// bytes: a byte order mark (FF FE), or a letter followed by a zero byte.
func readableOutput(output io.Reader) io.Reader {
	buffered := bufio.NewReader(output)
	start, _ := buffered.Peek(2) // a read error shows up again on the next read
	if len(start) == 2 && (start[0] == 0xFF && start[1] == 0xFE || start[0] != 0 && start[1] == 0) {
		return &utf16Reader{source: buffered}
	}
	return buffered
}

// utf16Reader turns UTF-16 (little-endian, as Windows writes it) into UTF-8.
type utf16Reader struct {
	source io.Reader
	carry  []byte // a last odd byte, kept until its partner arrives
	ready  []byte // UTF-8 text not yet returned
	err    error  // the source's error, returned once ready is empty
}

func (r *utf16Reader) Read(p []byte) (int, error) {
	chunk := make([]byte, 4096)
	for len(r.ready) == 0 && r.err == nil {
		n, err := r.source.Read(chunk)
		r.err = err
		r.decode(append(r.carry, chunk[:n]...))
	}
	if len(r.ready) == 0 {
		return 0, r.err
	}
	n := copy(p, r.ready)
	r.ready = r.ready[n:]
	return n, nil
}

// decode converts pairs of bytes to UTF-8 and keeps an odd last byte for the
// next read. A character outside the basic range (such as an emoji) that
// arrives split across two reads becomes a replacement mark. Windows tools do
// not print those characters.
func (r *utf16Reader) decode(data []byte) {
	units := make([]uint16, 0, len(data)/2)
	for len(data) >= 2 {
		units = append(units, binary.LittleEndian.Uint16(data))
		data = data[2:]
	}
	r.carry = append([]byte(nil), data...)
	for _, char := range utf16.Decode(units) {
		if char != byteOrderMark {
			r.ready = utf8.AppendRune(r.ready, char)
		}
	}
}
