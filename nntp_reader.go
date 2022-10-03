package nntp

import (
	"io"
	"net/textproto"
)

type rawDotReader struct {
	r     *textproto.Reader
	state int
}

// Read satisfies reads by streaming data up to and including the <DOT><CR><LF> termination sequence.
func (d *rawDotReader) Read(b []byte) (n int, err error) {
	// Run data through a simple state machine to
	// detect ending .\r\n line.
	const (
		stateBeginLine = iota // beginning of line; initial state; must be zero
		stateDot              // read . at beginning of line
		stateDotCR            // read .\r at beginning of line
		stateCR               // read \r (possibly at end of line)
		stateData             // reading data in middle of line
		stateEOF              // reached .\r\n end marker line
	)
	var c byte
	br := d.r.R
	bLen := len(b)
	for ; n < bLen && d.state != stateEOF; n++ {
		if c, err = br.ReadByte(); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			break
		}
		b[n] = c
		switch d.state {
		case stateBeginLine:
			if c == '.' {
				d.state = stateDot
			} else if c == '\r' {
				d.state = stateCR
			} else {
				d.state = stateData
			}

		case stateDot:
			if c == '\r' {
				d.state = stateDotCR
			} else if c == '\n' {
				d.state = stateEOF
				continue
			} else {
				d.state = stateData
			}

		case stateDotCR:
			if c == '\n' {
				d.state = stateEOF
				continue
			} else {
				d.state = stateData
			}

		case stateCR:
			if c == '\n' {
				d.state = stateBeginLine
			} else {
				d.state = stateData
			}

		case stateData:
			if c == '\r' {
				d.state = stateCR
			} else if c == '\n' {
				d.state = stateBeginLine
			}
		}
	}
	if err == nil && d.state == stateEOF {
		err = io.EOF
	}
	return
}
