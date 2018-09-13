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
$ ejson test.ejson
```

Would result in the following output:

```
export SECRET_SHELL_VARIABLE=<decrypted data>
```

You can then have your shell evaluate this output:

```shell
$ eval $(ejson test.ejson)
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

### Releasing

The release process is somewhat awkward right now. `ejson2env` is released in three ways:

* rubygem;
* `.deb` package; and
* homebrew formula.

Before releasing a new version, bump `/VERSION`, then run `make`, and commit the changes. Tag this
commit using `git tag vx.y.z`, e.g. `v1.0.0`.

In order to release the rubygem, find someone in the owners list at
https://rubygems.org/gems/ejson2env and ask them to add you, then:

1. `make`
2. `gem push pkg/ejson2env-x.y.z.gem`

To release the `.deb` package, edit the github release for the tag, and drop the `pkg/*.deb` in.

Releasing the homebrew package is more awkward. There is surely a more efficient way to do this but
my process is:

1. `dev clone homebrew-shopify`
1. Edit `ejson2env.rb`, changing the URL to reflect the new version. Also remove the `bottle do`
   paragraph.
1. Run `brew install ./ejson2env.rb` and change the `sha256` line to the SHA that is printed as an
   error.
1. `cd /usr/local/Homebrew/Library/Taps/shopify/homebrew-shopify`
1. `cp $(dev project-path homebrew-shopify)/ejson2env.rb .`
1. `brew install --build-bottle ejson2env`
1. `brew bottle ejson2env`
1. This will generate a file called `ejson2env--vx.y.z.high_sierra.tar.gz`. Rename it, turning `--`
   into `-`, e.g. `ejson2env-vx.y.z.high_sierra.tar.gz`.
1. Find some public place to upload this file. I use a personal S3 bucket but there are definitely
   better ways.
1. reset your changes in the homebrew tap and `dev cd homeshop`
1. Copy the `bottle do` paragraph printed by the `brew bottle` command into `ejson2env.rb`, and add
   `root_url "https://...`. The final expected URL is going to be `"#{root_url}/ejson2env...tar.gz`.
1. Commit and push directly to master.
1. `brew update && brew uninstall ejson2env && brew install ejson2env`. If this didn't work, or
   didn't correctly install from the bottle (i.e. took more than 10 seconds to install),
   troubleshoot, or revert.
