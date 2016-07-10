## Requirements

* Go 1.6

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

The executable is fully self-contained and can be distributed freely.

During development, you can define `FAKECLOUDGOV_DEBUG=yup` to make
the server fetch data files from the `data` directory instead of using
the files embedded into the executable at build time.

## Running Tests

```
go test
```

## Running the Example Client

An example OAuth2 client is in the `example-client` directory. To use
it, run:

```
cd example-client
npm install
npm start
```

Then visit http://localhost:8000/.

Note that the server, `./fake-cloud.gov`, must also be running in order
for the client to work.
