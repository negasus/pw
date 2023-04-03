package pw

import (
	"os"
	"testing"
)

func TestName1(t *testing.T) {
	tbl := NewTable()

	tbl.AddLine("fooqdwesdw", "bar", "baz", "", "12")
	tbl.AddLine("foo", "barqdswsqw qdsw", "baz")
	tbl.AddLine("foo", "bar", "bazqwds", "1 qwds")
	tbl.SetHeader("1", "2", "3", "4", "5", "6")

	tbl.Render(os.Stdout)
}
