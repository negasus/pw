package pw

import (
	"bytes"
	"fmt"
	"testing"
)

func printResult(s []byte) {
	for _, c := range s {
		fmt.Printf("%3d ", c)
	}
	fmt.Printf("\n")
	for _, c := range s {
		v := string(c)
		if c == '\n' {
			v = "\\n"
		}

		fmt.Printf("%3s ", v)
	}
	fmt.Printf("\n")
}

func TestName(t *testing.T) {
	s := `{"foo": "bar",
"num": 12,"num2": 12.34,
"b1": true,"b2": false,
"n1": null,
"o1": {"foo": "bar"},
"a1": [1, 2, 3, "foo", false, null, -15.33]}
`
	//http.ListenAndServe("127.0.0.1:2000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.Write([]byte(s))
	//}))
	//s = `{"a":{"b":1},"b":2}`

	cj := New([]byte(s))
	//cj.SetTab([]byte("  "))
	//cj.plainArray = true
	//err := PrintString(s)
	err := cj.Render()
	fmt.Printf("Error: %v\n", err)
}

func TestSimple(t *testing.T) {
	s := `{"foo":"bar"}`

	buf := bytes.NewBuffer(nil)
	err := FPrint(buf, []byte(s))
	if err != nil {
		t.Fatal(err)
	}

	expect := "{\n    \"foo\": " + string(ColorGreen) + "\"bar\"" + string(ColorReset) + "\n}\n"

	if buf.String() != expect {
		fmt.Printf("Result: %v\n", buf.Bytes())
		fmt.Printf("Expect: %v\n", []byte(expect))
		t.Fatal("bad result")
	}
}

func TestEscapedString(t *testing.T) {
	s := `{"foo":"bar \" baz"}`

	buf := bytes.NewBuffer(nil)
	err := FPrint(buf, []byte(s))
	if err != nil {
		t.Fatal(err)
	}

	expect := "{\n    \"foo\": " + string(ColorGreen) + "\"bar \\\" baz\"" + string(ColorReset) + "\n}\n"

	if buf.String() != expect {
		printResult(buf.Bytes())
		printResult([]byte(expect))
		t.Fatal("bad result")
	}
}
