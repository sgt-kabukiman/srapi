package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func main() {
	typeFlag := flag.String("type", "", "main type name, e.g. Game")
	typePluralFlag := flag.String("plural", "", "main type name as a plural, e.g. Games")

	flag.Parse()

	filename := strings.ToLower(*typeFlag) + "_collection.go"
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		panic(err)
	}

	t := *typeFlag
	tl := strings.ToLower(string(t[0])) + t[1:]
	tp := *typePluralFlag
	tpl := strings.ToLower(string(tp[0])) + tp[1:]
	s := string(tl[0])

	cTemplate, _ := template.ParseFiles("generate/collection.got")
	cTemplate.Execute(fp, map[string]string{
		"Sign":            s,
		"Type":            t,
		"TypeLower":       tl,
		"TypePlural":      tp,
		"TypePluralLower": tpl,
	})

	fp.Close()
	fmt.Printf("\n")
}
