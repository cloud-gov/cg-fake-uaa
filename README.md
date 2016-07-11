This is a fake User Account and Authentication ([UAA][]) server for
cloud.gov, useful for development and debugging.

## Build Requirements

* Go 1.6

Once built, the executable binary is fully self-contained and can be
distributed freely.

## Quick Start

First, get dependencies:

```
go get -d ./...
go get -u github.com/jteeuwen/go-bindata/...
```

Then generate and build:

```
go generate
go build
```

Finally, run the server:

```
./fake-cloud.gov
```

During development, you can define `FAKECLOUDGOV_DEBUG=yup` to make
the server fetch data files from the `data` directory instead of using
the files embedded into the executable at build time.

To learn about changing any of the runtime options, run
`./fake-cloud.gov -help`.

## Running Tests

```
go test
```

## Running the Example Client

A node-based example OAuth2 client is in the `example-client` directory.
To use it, run:

```
cd example-client
npm install
npm start
```

Then visit http://localhost:8000/.

Note that the server, `./fake-cloud.gov`, must also be running in order
for the client to work.

## Limitations

The fake server currently has a lot of limitations, most notably:

* The server has no support for refresh tokens.
* Only the [`openid` scope][] is supported. That is, the server is
  only really built for giving you the logged-in user's email
  address.

[UAA]: https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst
[`openid` scope]: https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst#scopes-authorized-by-the-uaa
