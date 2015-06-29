# net

[![net](https://godoc.org/github.com/mistifyio/util/net?status.png)](https://godoc.org/github.com/mistifyio/util/net)

Package net provides additional network utility functions.

## Usage

#### func  SplitHostPort

```go
func SplitHostPort(hostport string) (host string, port string, err error)
```
SplitHostPort splits a network address of the form "host", "host:port",
"[host]", "[host]:port", "[ipv6-host%zone]", or "[ipv6-host%zone]:port" into
host or ipv6-host%zone and port. Port will be an empty string if not supplied.

--
*Generated with [godocdown](https://github.com/robertkrimen/godocdown)*
