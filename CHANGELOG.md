# 2.0.9

* Security: keys are re-validated as POSIX shell names when exported, so `--trim-underscore` can no longer turn a key such as `_` into an empty name. In quiet mode an empty name produced a line beginning with `=`, which zsh runs as a command even under `eval "$(...)"`.
* Security: the documented shell invocation now quotes the command substitution, `eval "$(ejson2env ...)"`. Existing callers must update their invocation; upgrading the binary does not change caller shell code.
* Values containing tabs or beginning with `=` are no longer altered.

# 2.0.8 

Adds sanitization to output during decryption.

# 2.0.7

Bumps golang.org/x/crypto from 0.17.0 to 0.31.0

# 2.0.5

Bump to go 1.18.1 and ejson 1.3.3

# 2.0.4

Move to goreleaser, update dependencies
