[![Made with Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Go version number](https://img.shields.io/github/go-mod/go-version/kckrinke/go-cdk.svg)](https://github.com/kckrinke/go-cdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/kckrinke/go-cdk.svg)](https://pkg.go.dev/github.com/kckrinke/go-cdk)
[![GoReportCard](https://goreportcard.com/badge/github.com/kckrinke/go-cdk)](https://goreportcard.com/report/github.com/kckrinke/go-cdk)
[![Build status](https://github.com/kckrinke/go-cdk/workflows/codecov/badge.svg)](https://github.com/kckrinke/go-cdk/actions?query=workflow%3Acodecov)
[![codecov](https://codecov.io/gh/kckrinke/go-cdk/branch/trunk/graph/badge.svg?token=8AVBADVD1S)](https://codecov.io/gh/kckrinke/go-cdk)

# CDK - Curses Development Kit

This package provides the GDK equivalent for [CTK](https://github.com/kckrinke/go-ctk). This is not intended to be a parity of GDK in any way, rather this package simply fulfills the terminal drawing and basic event systems required by CTK.

Unless you're using CTK, you should really be using [TCell](https://github.com/gdamore/tcell) instead.

### Installing

```
go get -u github.com/kckrinke/go-cdk
```

### Building

A makefile has been included to assist in the development workflow.

```
> make help
usage: make {help|test|clean|demos}

  test: perform all available tests
  clean: cleans package  and built files
  demos: builds the boxes, mouse and unicode demos
```

## Example Usage

While CDK is not intended for direct usage, there are some simple demonstration applications provided.

### CDK Demo

A formal CDK application demonstrating the typical boilerplate setup.

* source code: [cdk-demo.go](_demos/cdk-demo.go)
* walkthrough: [pkg.go.dev](https://pkg.go.dev/github.com/kckrinke/go-cdk)

## Running the tests

CDK provides tests for color, event, runes and styles using the simulation screen. To run the tests, use the make-file for convenience:

```
> go test -v
testing cdk
  ... (per-test output, trimmed for brevity) ...
ok      github.com/kckrinke/go-cdk  0.171s
```

## Versioning

The current API is unstable and subject to change dramatically.

## License

This project is licensed under the Apache 2.0 license - see the [LICENSE.md](LICENSE.md) file for details.

## Authors and Contributors

* **Kevin C. Krinke** - *CDK author* - [kckrinke](https://github.com/kckrinke)

## Acknowledgments

* Thanks to [TCell](https://github.com/gdamore/tcell) for providing a solid and robust platform to build upon.

### TCell Authors and Contributors

* **Garrett D'Amore** - *Original author* - [gdamore](https://github.com/gdamore)
* **Zachary Yedidia** - *Contributor* - [zyedidia](https://github.com/zyedidia)
* **Junegunn Choi** - *Contributor* - [junegunn](https://github.com/junegunn)
* **Staysail Systems, Inc.** - *Support Provider* - [website](http://staysail.tech/)
