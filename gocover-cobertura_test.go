package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const SaveTestResults = false

type dirInfo struct {
	PkgPath string
}

func TestMain(t *testing.T) {
	fname := filepath.Join(os.TempDir(), "stdout")
	temp, _ := os.Create(fname)
	os.Stdout = temp
	main()
	outputBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fail()
	}
	outputString := string(outputBytes)
	if !strings.Contains(outputString, xml.Header) {
		t.Fail()
	}
	if !strings.Contains(outputString, coberturaDTDDecl) {
		t.Fail()
	}
}

func TestConvertParseProfilesError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || r != "Can't parse profiles" {
			t.Errorf("The code did not panic as expected; r = %+v", r)
		}
	}()

	pipe2rd, pipe2wr := io.Pipe()
	defer func() { pipe2rd.Close(); pipe2wr.Close() }()
	convert(strings.NewReader("invalid data"), pipe2wr)
}

func TestConvertOutputError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || r.(error).Error() != "io: read/write on closed pipe" {
			t.Errorf("The code did not panic as expected; r = %+v", r)
		}
	}()

	pipe2rd, pipe2wr := io.Pipe()
	pipe2wr.Close()
	defer func() { pipe2rd.Close() }()
	convert(strings.NewReader("mode: set"), pipe2wr)
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

func TestParseProfileDoesntExist(t *testing.T) {
	v := Coverage{}
	profile := Profile{FileName: "does-not-exist"}
	err := v.parseProfile(&profile)
	if err == nil || !strings.Contains(err.Error(), `can't find "does-not-exist"`) {
		t.Fatalf("Expected \"can't find\" error; got: %+v", err)
	}
}

func TestParseProfileNotReadable(t *testing.T) {
	v := Coverage{}
	profile := Profile{FileName: os.DevNull}
	err := v.parseProfile(&profile)
	if err == nil || !strings.Contains(err.Error(), `expected 'package', found 'EOF'`) {
		t.Fatalf("Expected \"expected 'package', found 'EOF'\" error; got: %+v", err)
	}
}

func TestParseProfilePermissionDenied(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "not-readable")
	defer os.Remove(tmpfile.Name())
	tmpfile.Chmod(000)
	v := Coverage{}
	profile := Profile{FileName: tmpfile.Name()}
	err = v.parseProfile(&profile)
	if err == nil || !strings.Contains(err.Error(), `permission denied`) {
		t.Fatalf("Expected \"permission denied\" error; got: %+v", err)
	}
}

func TestConvertSetMode(t *testing.T) {
	pipe1rd, err := os.Open("testdata/testdata_set.txt")
	if err != nil {
		t.Fatal("Can't parse testdata.")
	}

	pipe2rd, pipe2wr := io.Pipe()

	var convwr io.Writer = pipe2wr
	if SaveTestResults {
		testwr, err := os.Create("testdata/testdata_set.xml")
		if err != nil {
			t.Fatal("Can't open output testdata.", err)
		}
		defer testwr.Close()
		convwr = io.MultiWriter(convwr, testwr)
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
	if strings.TrimRight(p.Name, "/") != "./testdata" {
		t.Fatal(p.Name)
	}
	if p.Classes == nil || len(p.Classes) != 2 {
		t.Fatal()
	}

	c := p.Classes[0]
	if c.Name != "-" {
		t.Error()
	}
	if c.Filename != "./testdata/func1.go" {
		t.Errorf("Expected %s but %s", "./testdata/func1.go", c.Filename)
	}
	if c.Methods == nil || len(c.Methods) != 1 {
		t.Fatal()
	}
	if c.Lines == nil || len(c.Lines) != 4 {
		t.Errorf("Expected 4 lines but got %d", len(c.Lines))
	}

	m := c.Methods[0]
	if m.Name != "Func1" {
		t.Error()
	}
	if c.Lines == nil || len(c.Lines) != 4 {
		t.Errorf("Expected 4 lines but got %d", len(c.Lines))
	}

	var l *Line
	if l = m.Lines[0]; l.Number != 4 || l.Hits != 1 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}
	if l = m.Lines[1]; l.Number != 5 || l.Hits != 0 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}
	if l = m.Lines[2]; l.Number != 6 || l.Hits != 0 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}
	if l = m.Lines[3]; l.Number != 7 || l.Hits != 0 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}

	if l = c.Lines[0]; l.Number != 4 || l.Hits != 1 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}
	if l = c.Lines[1]; l.Number != 5 || l.Hits != 0 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}
	if l = c.Lines[2]; l.Number != 6 || l.Hits != 0 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}
	if l = c.Lines[3]; l.Number != 7 || l.Hits != 0 {
		t.Errorf("unmatched line: Number:%d, Hits:%d", l.Number, l.Hits)
	}

	c = p.Classes[1]
	if c.Name != "Type1" {
		t.Error()
	}
	if c.Filename != "./testdata/func2.go" {
		t.Errorf("Expected %s but %s", "./testdata/func2.go", c.Filename)
	}
	if c.Methods == nil || len(c.Methods) != 3 {
		t.Fatal()
	}
}
