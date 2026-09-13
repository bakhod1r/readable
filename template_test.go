package readable_test

import (
	htmltemplate "html/template"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/bakhod1r/readable"
)

func TestFuncMap(t *testing.T) {
	data := map[string]any{"Size": uint64(1536), "Elapsed": 3*time.Hour + 25*time.Minute, "Amount": int64(150050)}
	const src = `{{ bytes .Size }} · {{ duration .Elapsed }} · {{ money .Amount "USD" }} · {{ plural 2 "city" }}`
	const want = "1.5 KB · 3h 25m · 1,500.50 USD · cities"

	var b strings.Builder
	if err := template.Must(template.New("").Funcs(readable.FuncMap()).Parse(src)).Execute(&b, data); err != nil {
		t.Fatal(err)
	}
	if b.String() != want {
		t.Errorf("text/template = %q; want %q", b.String(), want)
	}

	b.Reset()
	h := htmltemplate.Must(htmltemplate.New("").Funcs(readable.FuncMap()).Parse(`<p>{{ maskEmail . }}</p>`))
	if err := h.Execute(&b, "john.doe@gmail.com"); err != nil {
		t.Fatal(err)
	}
	if b.String() != "<p>j***@gmail.com</p>" {
		t.Errorf("html/template = %q", b.String())
	}
}

func ExampleFuncMap() {
	tmpl := template.Must(template.New("").Funcs(readable.FuncMap()).Parse(
		`{{ bytes .Size }} in {{ duration .Took }}` + "\n"))
	_ = tmpl.Execute(os.Stdout, map[string]any{"Size": uint64(1 << 20), "Took": 1500 * time.Millisecond})
	// Output: 1 MB in 1s
}
