package mask

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type MaskedWriter struct {
	mask         string
	replacements []*regexp.Regexp
	source       io.Reader
	target       io.Writer
}

func (w *MaskedWriter) Write() {
	scanner := bufio.NewScanner(w.source)

	for scanner.Scan() {
		fmt.Fprintln(w.target, maskLine(scanner.Text(), w.mask, w.replacements))
	}
}

func NewMaskedWriter(r *Masks, s io.Reader, t io.Writer) *MaskedWriter {
	replacements := r.Compile()

	return &MaskedWriter{
		mask:         r.MaskChar,
		replacements: replacements,
		source:       s,
		target:       t,
	}
}

func maskLine(line, replace string, re []*regexp.Regexp) string {
	for _, mask := range re {
		line = mask.ReplaceAllString(line, replace)
	}
	return line
}
