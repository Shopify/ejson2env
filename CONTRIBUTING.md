# Contributing

Thank you for considering contributing to ejson2env!

## Getting Started

- Review this document and the [Code of Conduct](CODE_OF_CONDUCT.md).

- Setup a [Go development environment](https://golang.org/doc/install#install)
if you haven't already.

- Get the Shopify version the project by using Go get:

```shell
$ go get -u github.com/Shopify/ejson2env/cmd/ejson2env
```

- Fork this project on GitHub. :octocat:

- Setup your fork as a remote for your project:

```
$ cd $GOPATH/src/github.com/Shopify/ejson2env
$ git remote add <your username> <your fork's remote path>
```

## Work on your feature

- Create your feature branch based off of the `master` branch. (It might be
worth doing a `git pull` if you haven't done one in a while.)

```
$ git checkout master
$ git pull
$ git checkout -b <the name of your branch>
```

- Code/write! :keyboard:

    - If working on code, please run `go fmt` and `golint` while you work on
your change, to clean up your formatting/check for issues.

- Push your changes to your fork's remote:

```
$ git push -u <your username> <the name of your branch>
```

## Send in your changes

- Sign the [Contributor License Agreement](https://cla.shopify.com).

- Open a PR against Shopify/ejson2env!

## Releasing

The release process is somewhat awkward right now. `ejson2env` is released in
four ways:

* `linux-amd64` and `darwin-amd64`
* rubygem;
* `.deb` package; and
* homebrew formula

Before releasing a new version, bump `/VERSION`, then run `make`, and commit
the changes. Tag this commit using `git tag vx.y.z`, e.g. `v1.0.0`.

In order to release the rubygem, find someone in the owners list at
https://rubygems.org/gems/ejson2env and ask them to add you, then:

1. `make`
2. `gem push pkg/ejson2env-x.y.z.gem`

To release the `.deb` package, edit the github release for the tag, and drop
the `pkg/*.deb` in.

Releasing the homebrew package is more awkward. There is surely a more
efficient way to do this but current process is:

1. `dev clone homebrew-shopify`
1. Edit `ejson2env.rb`, changing the URL to reflect the new version. Also
   remove the `bottle do` paragraph.
1. Run `brew install ./ejson2env.rb` and change the `sha256` line to the SHA
   that is printed as an error.
1. `cd /usr/local/Homebrew/Library/Taps/shopify/homebrew-shopify`
1. `cp $(dev project-path homebrew-shopify)/ejson2env.rb .`
1. `brew install --build-bottle ejson2env`
1. `brew bottle ejson2env`
1. This will generate a file called `ejson2env--vx.y.z.high_sierra.tar.gz`.
   Rename it, turning `--` into `-`, e.g. `ejson2env-vx.y.z.high_sierra.tar.gz`.
1. Find some public place to upload this file. I use a personal S3 bucket but
   there are definitely better ways.
1. reset your changes in the homebrew tap and `dev cd homeshop`
1. Copy the `bottle do` paragraph printed by the `brew bottle` command into
   `ejson2env.rb`, and add `root_url "https://...`. The final expected URL is
    going to be `"#{root_url}/ejson2env...tar.gz`.
1. Commit and push directly to master.
1. `brew update && brew uninstall ejson2env && brew install ejson2env`. If this
   didn't work, or didn't correctly install from the bottle (i.e. took more
   than 10 seconds to install), troubleshoot, or revert.
