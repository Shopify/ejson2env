# ejson2env

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
