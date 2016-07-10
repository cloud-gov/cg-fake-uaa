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
