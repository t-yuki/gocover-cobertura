[![Build Status](https://travis-ci.org/t-yuki/gocover-cobertura.svg?branch=master)](https://travis-ci.org/t-yuki/gocover-cobertura)
[![Coverage Status](https://coveralls.io/repos/github/t-yuki/gocover-cobertura/badge.svg?branch=master)](https://coveralls.io/github/t-yuki/gocover-cobertura?branch=master)

go tool cover XML (Cobertura) export
====================================

This is a simple helper tool for generating XML output in [Cobertura](http://cobertura.sourceforge.net/) format
for CIs like [Jenkins](https://wiki.jenkins-ci.org/display/JENKINS/Cobertura+Plugin) and others
from [go tool cover](https://code.google.com/p/go.tools/) output.

Installation
------------

Just type the following to install the program and its dependencies:

    $ go get code.google.com/p/go.tools/cmd/cover
    $ go get github.com/t-yuki/gocover-cobertura

Usage
-----

`gocover-cobertura` reads from the standard input:

    $ go test -coverprofile=coverage.txt -covermode count github.com/gorilla/mux
    $ gocover-cobertura < coverage.txt > coverage.xml

Some flags can be passed (each flag should only be used once):

- `-ignore-dirs PATTERN`

  ignore directories matching `PATTERN` regular expression. Full
  directory names are matched, as
  `github.com/t-yuki/gocover-cobertura` (and so `github.com/t-yuki`
  and `github.com`), examples of use:
  ```
  # A specific directory
  -ignore-dirs '^github\.com/t-yuki/gocover-cobertura/testdata$'
  # All directories autogen and any of their subdirs
  -ignore-dirs '/autogen$'
  ```

- `-ignore-files PATTERN`

  ignore files matching `PATTERN` regular expression. Full file names
  are matched, as `github.com/t-yuki/gocover-cobertura/profile.go`,
  examples of use:
  ```
  # A specific file
  -ignore-files '^github\.com/t-yuki/gocover-cobertura/profile\.go$'
  # All files ending with _gen.go
  -ignore-files '_gen\.go$'
  # All files in a directory autogen (or any of its subdirs)
  -ignore-files '/autogen/'
  ```

- `-ignore-gen-files`

  ignore generated files. Typically files containing a comment
  indicating that the file has been automatically generated. See
  `genCodeRe` regexp in [ignore.go](ignore.go).

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)

Thanks
------

This tool is originated from [gocov-xml](https://github.com/AlekSi/gocov-xml) by [Alexey Palazhchenko (AlekSi)](https://github.com/AlekSi)
