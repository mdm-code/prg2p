<h1 align="center">
  <div>
    <img src="https://raw.githubusercontent.com/mdm-code/mdm-code.github.io/main/prg2p_logo.png" alt="logo"/>
  </div>
</h1>

<h4 align="center">Grapheme-to-phoneme converter for Polish in Go</h4>

<div align="center">
<p>
    <a href="https://github.com/mdm-code/prg2p/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/mdm-code/prg2p/workflows/CI/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/prg2p">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/prg2p/branch/main/graphs/badge.svg?branch=main">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/prg2p">
    </a>
    <a href="https://goreportcard.com/report/github.com/mdm-code/prg2p">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/prg2p">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/prg2p">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
</p>
</div>

The `prg2p` package implements a grapheme-to-phoneme rule based converter for
Polish.

It provides a standalone command-line program to process data efficiently on
the terminal and exposes the public API components of the package for code
reuse.

Consult the [package documentation](https://pkg.go.dev/github.com/mdm-code/prg2p)
or check [Usage](#usage) section below to see how to use `prg2p` in your code.


## Installation

To add package to a Go project dependencies run the following command:

```sh
go get github.com/mdm-code/prg2p
```

In order to use the CLI program, you need to use this command:

```sh
go install github.com/mdm-code/prg2p@latest
```

Here, you can use the `@latest` or any version you find appropriate for that
matter.


## Usage

Type `prg2p -h` from the terminal after installing executables as described
[here](#installation) to see how to use `prg2p` command-line interface.

Here is how you can use the public API of the `prg2p` package in your code:

```go
package main

import (
	"fmt"

	"github.com/mdm-code/prg2p"
)

func main() {
	r := prg2p.Rules()
	g2p, err := prg2p.Load(r)

	// Iterate over words to get their phonemic transcripts
	var trans []string
	for _, w := range []string{"ala", "ma", "kota"} {
		t, err := g2p.Transcribe(w, false)
		if err != nil {
			fmt.Println(err)
			continue
		}
		trans = append(trans, t...)
	}
	for _, t := range trans {
		fmt.Println(t)
	}
}
```


## Development

All necessary development tools are in the `Makefile`. Calling `make test`
consecutively invokes `go fmt`, `go vet`, `golint` and `go test`. CI/CD is
handled by Github workflows. Remember to install `golint` before testing and
building:

```sh
go install golang.org/x/lint/golint@latest
```

Happy coding!


## License

Copyright (c) 2022 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
