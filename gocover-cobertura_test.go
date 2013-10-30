package main

import (
	"encoding/xml"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

type dirInfo struct {
	PkgPath string
}

func TestConvertEmpty(t *testing.T) {
	data := `mode: set`

	pipe2rd, pipe2wr := io.Pipe()
	go convert(strings.NewReader(data), pipe2wr)

	v := Coverage{}
	dec := xml.NewDecoder(pipe2rd)
	dec.Decode(&v)

	if v.XMLName.Local != "coverage" {
		t.Error()
	}
	if v.Sources == nil {
		t.Fatal()
	}
	if v.Packages != nil {
		t.Fatal()
	}
}

func TestConvertFunc1SetMode(t *testing.T) {
	tmpl, err := template.ParseFiles("testdata/func1_test_set.txt")
	if err != nil {
		t.Fatal("Can't parse testdata.")
	}
	dirInfo := dirInfo{}
	dirInfo.PkgPath = reflect.TypeOf(Coverage{}).PkgPath()

	pipe1rd, pipe1wr := io.Pipe()
	go func() {
		err := tmpl.Execute(pipe1wr, dirInfo)
		if err != nil {
			t.Error("Can't execute template.")
			panic("tmpl.Execute failed")
		}
		pipe1wr.Close()
	}()

	pipe2rd, pipe2wr := io.Pipe()

	var convwr io.Writer = pipe2wr
	testwr, err := os.Create("testdata/func1_test_set.xml")
	if err == nil {
		convwr = io.MultiWriter(convwr, testwr)
	} else {
		t.Log("Can't open output testdata. ignoring...")
	}

	go convert(pipe1rd, convwr)

	v := Coverage{}
	dec := xml.NewDecoder(pipe2rd)
	dec.Decode(&v)

	if v.XMLName.Local != "coverage" {
		t.Error()
	}

	if v.Sources == nil {
		t.Fatal()
	}

	if v.Packages == nil || len(v.Packages) != 1 {
		t.Fatal()
	}

	p := v.Packages[0]
	if p.Name != dirInfo.PkgPath+"/testdata" {
		t.Fatal()
	}
	if p.Classes == nil || len(p.Classes) != 1 {
		t.Fatal()
	}

	c := p.Classes[0]
	if c.Name != "-" {
		t.Error()
	}
	if c.Filename != dirInfo.PkgPath+"/testdata/func1.go" {
		t.Errorf("Expected %s but %s", dirInfo.PkgPath+"/testdata/func1.go", c.Filename)
	}
	if c.Methods == nil || len(c.Methods) != 1 {
		t.Fatal()
	}
	if c.Lines == nil || len(c.Lines) != 2 {
		t.Fatal()
	}

	m := c.Methods[0]
	if m.Name != "Func1" {
		t.Error()
	}
	if m.Lines == nil || len(m.Lines) != 2 {
		t.Fatal()
	}

	l1 := m.Lines[0]
	if l1.Number != 5 {
		t.Error()
	}
	if l1.Hits != 1 {
		t.Error()
	}

	l2 := m.Lines[1]
	if l2.Number != 6 {
		t.Error()
	}
	if l2.Hits != 0 {
		t.Error()
	}

	l1 = c.Lines[0]
	if l1.Number != 5 {
		t.Error()
	}
	if l1.Hits != 1 {
		t.Error()
	}

	l2 = c.Lines[1]
	if l2.Number != 6 {
		t.Error()
	}
}
