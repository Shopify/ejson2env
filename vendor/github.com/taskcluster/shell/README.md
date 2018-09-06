# shell

[![GoDoc](https://godoc.org/github.com/taskcluster/shell?status.svg)](https://godoc.org/github.com/taskcluster/shell)
[![Build Status](https://travis-ci.org/taskcluster/shell.svg?branch=master)](http://travis-ci.org/taskcluster/shell)
[![License](https://img.shields.io/badge/license-MPL%202.0-orange.svg)](https://github.com/taskcluster/shell/blob/master/LICENSE)

Escape command line parameters to be executed on the shell.

Ported from https://github.com/xxorax/node-shell-escape.

## Usage

```
package main

import (
	"fmt"

	"github.com/taskcluster/shell"
)

func main() {
	fmt.Println(
		shell.Escape(
			"curl",
			"-v",
			"-H",
			"Location;",
			"-H",
			"User-Agent: dave#10",
			"http://www.daveeddy.com/?name=dave&age=24",
		),
	)
}
```

produces...

```
curl -v -H 'Location;' -H 'User-Agent: dave#10' 'http://www.daveeddy.com/?name=dave&age=24'
```
