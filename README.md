[![codecov](https://codecov.io/gh/kckrinke/go-cdk/branch/trunk/graph/badge.svg?token=8AVBADVD1S)](https://codecov.io/gh/kckrinke/go-cdk)

# CDK - Curses Display Kit

This package provides the GDK equivalent for [CTK](https://github.com/kckrinke/go-ctk). This is not intended to be a parity of GDK in any way, rather this package simply fulfills the terminal drawing and basic event systems required by CTK.

Unless you're using CTK, you should really be using [tcell](https://github.com/gdamore/tcell) instead.

### Installing

```
go get -u github.com/kckrinke/go-cdk
```

### Building

A makefile has been included to assist in the development workflow.

```
$ make help
usage: make {help|test|clean|demos}

  test: perform all available tests
  clean: cleans package  and built files
  demos: builds the boxes, mouse and unicode demos
```

## Example Usage

While CDK is not intended for direct usage, there are some simple demonstration applications provided.

### Boxes

Simple demostration of rendering randomly sized and positioned rectangles on the screen. Pressing ESC will exit the application and inform you of the total number of boxes drawn and the average time for each box's rendering.

* See: [boxes.go](_demos/boxes.go)

### Mouse

Displays an information box and provides a drawing area for testing and demonstrating the mouse interaction events.

* See: [mouse.go](_demos/mouse.go)

### Unicode

Somewhat complicated code demonstration of working with advanced unicode characters. The actual application simply renders a bunch of non-ASCII characters of various widths and encoding.

* See: [unicode.go](_demos/unicode.go)

## Running the tests

CDK provides tests for color, event, runes and styles using the simulation screen. To run the tests, use the make-file for convenience:

```
$ make test
testing cdk
ok      github.com/kckrinke/go-cdk  0.171s
?       github.com/kckrinke/go-cdk/encoding [no test files]
ok      github.com/kckrinke/go-cdk/terminfo 0.052s
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
