package main

import (
	"fmt"
	"flag"
	"net/http"
	"html/template"
	_ "embed"
)

var (
	//go:embed "template.html"
	templateHTML string
)

func die(message string, err error) {
	fmt.Println(message)
	fmt.Println(err)
	panic("failed")
}

type (
	Pattern struct {
		Size int
		Pad string
	}
)

func main() {
	size := flag.Int("size", 25, "pattern size")
	pad := flag.Int("pad", 4, "id padding (with 0s)")
	bind := flag.String("bind", ":10987", "local binding to use for server")
	flag.Parse()
	template, err := template.New("t").Parse(templateHTML)
	if err != nil {
		die("unable to parse template", err)
	}
	padding := *pad
	if padding < 1 {
		die("invalid padding", fmt.Errorf("< 1"))
	}
	padString := ""
	for padding > 0 {
		padString = fmt.Sprintf("0%s", padString)
		padding = padding - 1
	}
	obj := Pattern{Size: *size + 1, Pad: padString}
	if obj.Size <= 0 {
		die("invalid size", fmt.Errorf("< 0"))
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := template.Execute(w, obj)
		if err != nil {
			fmt.Println(fmt.Sprintf("failed to execute template: %v", err))
		}
    })
	binding := *bind
	fmt.Println(binding)
	if err := http.ListenAndServe(binding, nil); err != nil {
		die("unable to bind", err)
	}
}
