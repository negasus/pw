package pw

import "fmt"

func (pw *PrettyWriter) RenderXML() error {
	if !pw.next() {
		return pw.error()
	}

	return pw.renderXML()
}

func (pw *PrettyWriter) renderXML() error {
	return fmt.Errorf("not implemented")
}
