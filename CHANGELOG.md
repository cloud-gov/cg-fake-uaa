# Change log

All notable changes to this project will be documented in this file,
which uses the format described in
[keepachangelog.com](http://keepachangelog.com/). This project adheres
to [Semantic Versioning](http://semver.org/).

## [Unreleased][unreleased]

## [1.1.0][] - 2017-04-20

* The `/oauth/token` endpoint now supports the `refresh_token` grant
  type, i.e. the server can now be used to refresh access tokens.

* The `-token-lifetime` command-line argument can be used to specify
  a lifetime for access tokens, e.g. `10m` or `10s`.

* The example JS client now refreshes the user's access token automatically.

## [1.0.4][] - 2017-03-23

Yet another attempt at rebuilding properly with `goxc`.

## [1.0.3][] - 2017-03-23

Another attempt at rebuilding properly with `goxc`.

## [1.0.2][] - 2017-03-23

No changes to the source--this is just a re-release of 1.0.1, but
including linux as a build target.

## [1.0.1][] - 2016-07-16

* Show version number on CLI and login page.

## 1.0.0 - 2016-07-16

Initial release.

[unreleased]: https://github.com/18F/cg-fake-uaa/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/18F/cg-fake-uaa/compare/v1.0.4...v1.1.0
[1.0.4]: https://github.com/18F/cg-fake-uaa/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/18F/cg-fake-uaa/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/18F/cg-fake-uaa/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/18F/cg-fake-uaa/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/18F/cg-fake-uaa/compare/v0.0.1...v1.0.0
