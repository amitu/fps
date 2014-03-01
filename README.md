# Flash Police Server

Written in `go`. Available as a executable as well as a library.

## Obtaining

Assuming you have `go` installed and `GOROOT` environment variable set.

```shell
$ go get github.com/amitu/fps
```

This will download `amitu/fps` from github and build/install it as a library.

## Executable

To build the executable:

```shell
$ go install github.com/amitu/fps/fps
```

This will build executable named `fps` and install it in your `$GOPATH/bin`. 

To run the executable:

```shell
$ $GOROOT/bin/fps policy.xml:localhost:843 policy2.xml:localhost:8000
```

This will launch two servers, one listening on `883`, serving content of `policy.xml` and the other listening on `8000` and serving the content of `policy2.xml`. You can pass as many `file:host:port` arguments as you want.

To serve the requests, `fps` will launch 10 workers, these 10 workers are shared with all servers started. You can pass `--workers=20` to change the number of workers. If you want to start different set of workers for different servers, consider using `fps` as a library.

`fps` by default runs in `non-strict` mode. The policy server protocol expects a request from client, before it sends the content of policy file. By default this implementation simply writes the content of policy file as soon as a client is connected, without reading anything from client. You can call `fps` with `--strict=true` to change this behaviour (*not yet implemented*).

To cleanly stop all servers, send SIGINT to the server, or press Ctrl-C. 

## `fps` Library

The library can be used by importing `github.com/amitu/pfs`. The documentation of library is available at [godocs](http://godoc.org/github.com/amitu/fps).

