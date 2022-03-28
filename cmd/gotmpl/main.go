package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/huandu/xstrings"
)

func stdinData() map[string]interface{} {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var data = map[string]interface{}{}
	err = json.Unmarshal(stdin, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

var output string
var input string

func main() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(),
			`Usage: gotmpl [-i data.json] [-o output.txt] [key=value] template [template ...]

Version 0.2.0
https://github.com/NateScarlet/gotmpl
`)
	}
	flag.StringVar(&output, "o", "", "output file path")
	flag.StringVar(&input, "i", "", "input json data file path, pass `-` for stdin")
	flag.Parse()

	data := map[string]interface{}{} // sprig dict functions require map[string]interface{}
	files := []string{}

	if output != "" {
		p, err := filepath.Abs(output)
		if err != nil {
			log.Fatal(err)
		}
		data["Name"] = strings.SplitN(filepath.Base(p), ".", 2)[0]
		data["Package"] = filepath.Base(filepath.Dir(p))
	}

	if input == "-" {
		for k, v := range stdinData() {
			data[k] = v
		}
	} else if input != "" {
		d, err := ioutil.ReadFile(input)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(d, &data)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, i := range flag.Args() {
		if strings.Contains(i, "=") {
			kv := strings.SplitN(i, "=", 2)
			if kv[0] == "" {
				files = append(files, kv[1])
				continue
			}
			data[kv[0]] = kv[1]
		} else {
			files = append(files, i)
		}
	}

	if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	funcs := sprig.TxtFuncMap()
	funcs["cwd"] = os.Getwd
	funcs["args"] = func() []string { return os.Args }
	funcs["templateFiles"] = func() []string { return files }
	funcs["__file__"] = func() string { return output }
	funcs["absPath"] = filepath.Abs
	funcs["upperFirst"] = xstrings.FirstRuneToUpper
	funcs["lowerFirst"] = xstrings.FirstRuneToLower

	tmpl, err := template.
		New(filepath.Base(files[0])).
		Funcs(funcs).
		ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, data)
	if err != nil {
		log.Fatal(err)
	}
	if output != "" {
		var err = ioutil.WriteFile(output, b.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Print(b.String())
	}
}
