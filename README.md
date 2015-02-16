lessc-watch
===========

A tool written in [Go](http://golang.org), which keeps track on file changes on [LESS](http://lesscss.org) stylesheets and compiles them automatically into CSS.

### Prerequisites

Make sure to have [Go installed and setup correctly](http://golang.org/doc/install) (I'm using Go 1.4.1 – no other versions are tested) and did the [$GOPATH thing](https://github.com/golang/go/wiki/GOPATH). LESS (especially `lessc`) has to be installed and made available in your $PATH as well (preferably via [`npm -g`](https://www.npmjs.com/package/less)).

### Install

To install `lessc-watch` via Go's package manager, try `go get github.com/jankassel/lessc-watch`. 

### Usage 

```
lessc-watch [-x] input-directory [output-directory]
```

The listener will watch for file changes in the `input-directory`. `output-directory` states where to put compiler output (i.e. the CSS). If omitted, all output will go into the same directory as `input-directory`. `-x` will enable compression, passing the same parameter to `lessc`.

`lessc-watch` will tell you when and which file has been compiled, and if any errors occured.

### Notes

As I'm working with [Sublime Text](http://www.sublimetext.com/3), I stumbled upon `atomic_save` when beginning to use `fsnotify`. `lessc-watch` currently just listens for Write events, which requires `atomic_save` [to be set to `false`](http://stackoverflow.com/a/20639093) in your preferences if you're using Sublime Text as well.
