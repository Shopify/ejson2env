[![Build Status](https://travis-ci.org/Shopify/ejson2env.svg?branch=master)](https://travis-ci.org/Shopify/ejson2env)
[![codecov](https://codecov.io/gh/Shopify/ejson2env/branch/master/graph/badge.svg)](https://codecov.io/gh/Shopify/ejson2env)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shopify/ejson2env)](https://goreportcard.com/report/github.com/Shopify/ejson2env)

# ejson2env

`ejson2env` is a tool to simplify storing secrets that should be accessible in the shell environment in your git repo. `ejson2env` is based on the [ejson library](https://github.com/Shopify/ejson) and extends the `ejson` file format.

`ejson2env` exports all of the values in the `environment` object in the `ejson` file to the shell environment.

For example, with the below `ejson` file:

```json
{
    "_public_key": "<public key here>",
    "environment": {
        "SECRET_SHELL_VARIABLE": "<encrypted data>"
    }
}
```

Running:

```shell
$ ejson2env test.ejson
```

Would result in the following output:

```
export SECRET_SHELL_VARIABLE=<decrypted data>
```

You can then have your shell evaluate this output:

```shell
$ eval $(ejson2env test.ejson)
```

## Using ejson2env

`ejson2env`'s usage information is described in it's included [manual page](/man/ejson2env.1.ronn).

## Installing ejson2env

`ejson2env` is available through a number of different routes and package managers. If you plan on modifying `ejson2env`, it is suggested that you install via `go get`.

### Go

`ejson2env` can be installed using the regular `go get` tool:

```shell
$ go get -u github.com/Shopify/ejson2env/cmd/ejson2env
```

You can then find the compiled binary in `$GOPATH/bin`

### Debian Package

You can download the latest version of the Debian package from [the releases page](https://github.com/Shopify/ejson2env/releases).

Install the downloaded package by calling:

```shell
$ dpkg -i ejson2env_1.0.3_amd64.deb
```

### RubyGems

You can install `ejson2env` using Ruby's Gem tool:

```shell
$ gem install ejson2env
```

### Homebrew

Provided your install of Homebrew is configured to pull from [Shopify's Homebrew repo](https://github.com/shopify/homebrew-shopify), you can install `ejson2env` by calling:

```shell
$ brew install ejson2env
```

## Contributing

Please review the [Contributing document](CONTRIBUTING.md) if you are
interested in helping improve ejson2env!
