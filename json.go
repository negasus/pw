package pw

import (
	"bytes"
	"fmt"
	"io"
)

func (pw *PrettyWriter) RenderJSON() error {
	if !pw.next() {
		return pw.error()
	}

	return pw.renderJSON()
}

func (pw *PrettyWriter) renderJSON() error {
	if !pw.printValue() && pw.err != io.EOF {
		return pw.error()
	}

	if !pw.writeByte('\n') {
		return pw.error()
	}

	return nil
}

func (pw *PrettyWriter) printValue() bool {
	switch pw.ch {
	case '"':
		return pw.withColor(pw.colorString, pw.printString)
	case '{':
		return pw.printObject()
	case '[':
		return pw.printArray()
	case 't', 'T':
		return pw.withColor(pw.colorBool, pw.printTrue)
	case 'f', 'F':
		return pw.withColor(pw.colorBool, pw.printFalse)
	case 'n', 'N':
		return pw.withColor(pw.colorNull, pw.printNull)
	}
	if pw.ch == '-' || (pw.ch >= '0' && pw.ch <= '9') {
		return pw.withColor(pw.colorNumber, pw.printNumber)
	}
	pw.err = fmt.Errorf("unexcpected char %s", string(pw.ch))
	return false
}

func (pw *PrettyWriter) writeByte(b ...byte) bool {
	_, err := pw.w.Write(b)
	if err != nil {
		pw.err = err
		return false
	}
	return true
}

func (pw *PrettyWriter) writeTabs() bool {
	for i := 0; i < pw.t; i++ {
		if !pw.writeByte(pw.tab...) {
			return false
		}
	}
	return true
}

func (pw *PrettyWriter) printObject() bool {
	if !pw.writeByte('{') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch == '}' {
		if !pw.next() {
			return false
		}
		return pw.writeByte('}')
	}
	if !pw.writeByte('\n') {
		return false
	}
	pw.t++
	if !pw.writeByte(bytes.Repeat(pw.tab, pw.t)...) {
		return false
	}
	for {
		if pw.ch != '"' {
			pw.err = fmt.Errorf("unexpected char %s, expect quote", string(pw.ch))
			return false
		}
		if !pw.printString() {
			return false
		}
		if pw.ch != ':' {
			pw.err = fmt.Errorf("unexpected char %s, expect colon", string(pw.ch))
			return false
		}
		if !pw.writeByte(':', ' ') {
			return false
		}
		if !pw.next() {
			return false
		}
		if !pw.printValue() {
			return false
		}
		if pw.ch == ',' {
			if !pw.writeByte(',', '\n') {
				return false
			}
			if !pw.writeByte(bytes.Repeat(pw.tab, pw.t)...) {
				return false
			}
			if !pw.next() {
				return false
			}
			continue
		}
		if pw.ch == '}' {
			if !pw.writeByte('\n') {
				return false
			}
			pw.t--
			if !pw.writeByte(bytes.Repeat(pw.tab, pw.t)...) {
				return false
			}
			if !pw.writeByte('}') {
				return false
			}

			return pw.next()
		}
		pw.err = fmt.Errorf("unexpected char %s, expect comma or right brace", string(pw.ch))
		return false
	}
}

func (pw *PrettyWriter) printString() bool {
	if !pw.writeByte('"') {
		return false
	}
	var escaped bool
	for {
		if !pw.nextNoSpace() {
			return false
		}
		if pw.ch == '\\' {
			escaped = !escaped
			if !pw.writeByte('\\') {
				return false
			}
			continue
		}
		if pw.ch == '"' && !escaped {
			if !pw.writeByte('"') {
				return false
			}
			return pw.next()
		}
		if !pw.writeByte(pw.ch) {
			return false
		}
		escaped = false
	}
}

func (pw *PrettyWriter) printArray() bool {
	if !pw.writeByte('[') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch == ']' {
		if !pw.next() {
			return false
		}
		return pw.writeByte(']')
	}
	if !pw.plainArray {
		if !pw.writeByte('\n') {
			return false
		}
		pw.t++
		if !pw.writeByte(bytes.Repeat(pw.tab, pw.t)...) {
			return false
		}
	}
	for {
		if !pw.printValue() {
			return false
		}
		if pw.ch == ',' {
			if !pw.writeByte(',') {
				return false
			}
			if !pw.plainArray {
				if !pw.writeByte('\n') {
					return false
				}
				if !pw.writeByte(bytes.Repeat(pw.tab, pw.t)...) {
					return false
				}
			} else {
				if !pw.writeByte(' ') {
					return false
				}
			}
			if !pw.next() {
				return false
			}
			continue
		}
		if pw.ch == ']' {
			if !pw.next() {
				return false
			}
			if !pw.plainArray {
				if !pw.writeByte('\n') {
					return false
				}
				pw.t--
				if !pw.writeByte(bytes.Repeat(pw.tab, pw.t)...) {
					return false
				}
			}
			return pw.writeByte(']')
		}
		pw.err = fmt.Errorf("unexpected char %s, expect comma or right bracket", string(pw.ch))
		return false
	}
}

func (pw *PrettyWriter) printTrue() bool {
	if !pw.writeByte('t') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'r' && pw.ch != 'R' {
		pw.err = fmt.Errorf("unexpected char %s, expect r", string(pw.ch))
		return false
	}
	if !pw.writeByte('r') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'u' && pw.ch != 'U' {
		pw.err = fmt.Errorf("unexpected char %s, expect u", string(pw.ch))
		return false
	}
	if !pw.writeByte('u') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'e' && pw.ch != 'E' {
		pw.err = fmt.Errorf("unexpected char %s, expect e", string(pw.ch))
		return false
	}
	if !pw.writeByte('e') {
		return false
	}
	return pw.next()
}

func (pw *PrettyWriter) printFalse() bool {
	if !pw.writeByte('f') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'a' && pw.ch != 'A' {
		pw.err = fmt.Errorf("unexpected char %s, expect a", string(pw.ch))
		return false
	}
	if !pw.writeByte('a') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'l' && pw.ch != 'L' {
		pw.err = fmt.Errorf("unexpected char %s, expect l", string(pw.ch))
		return false
	}
	if !pw.writeByte('l') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 's' && pw.ch != 'S' {
		pw.err = fmt.Errorf("unexpected char %s, expect s", string(pw.ch))
		return false
	}
	if !pw.writeByte('s') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'e' && pw.ch != 'E' {
		pw.err = fmt.Errorf("unexpected char %s, expect e", string(pw.ch))
		return false
	}
	if !pw.writeByte('e') {
		return false
	}
	return pw.next()
}

func (pw *PrettyWriter) printNull() bool {
	if !pw.writeByte('n') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'u' && pw.ch != 'U' {
		pw.err = fmt.Errorf("unexpected char %s, expect u", string(pw.ch))
		return false
	}
	if !pw.writeByte('u') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'l' && pw.ch != 'L' {
		pw.err = fmt.Errorf("unexpected char %s, expect l", string(pw.ch))
		return false
	}
	if !pw.writeByte('l') {
		return false
	}
	if !pw.next() {
		return false
	}
	if pw.ch != 'l' && pw.ch != 'L' {
		pw.err = fmt.Errorf("unexpected char %s, expect l", string(pw.ch))
		return false
	}
	if !pw.writeByte('l') {
		return false
	}
	return pw.next()
}

// todo: check on '-.10' '.10', exp
func (pw *PrettyWriter) printNumber() bool {
	if !pw.writeByte(pw.ch) {
		return false
	}
	var hasDot bool
	for {
		if !pw.next() {
			return false
		}
		if pw.ch == '.' {
			if hasDot {
				pw.err = fmt.Errorf("unexpected char %s, expect number", string(pw.ch))
				return false
			}
			hasDot = true
			if !pw.writeByte('.') {
				return false
			}
			continue
		}
		if !isDigit(pw.ch) {
			return true
		}
		if !pw.writeByte(pw.ch) {
			return false
		}
	}
}
