# Fabulously Easy Reverse Proxy

Super easy to use reverse proxy with simple to understand configuration.

## Why?

I run a small personal website and I found nginx way too annoying to set up
with many backend hosts & redirects.

### Goals

#### Firm (Must)

- Configure self-descriptively with a modern markup language
- Fail descriptively when something is misconfigured
- Address the use case of a small personal portfolio website

#### Soft (Nice)

- Handle headers according to reverse proxy "standards"

### Non Goals

- Be feature packed
    - Address only those things that are required most often by the main use case
- Serve with blazing speed 
    - Go's HTTP performance is pretty good
    - Intended for hobby projects with small audiences

## Usage

```go
go build && ./ferp --configuration-file /path/to/configuration.yaml
```

## Configuration

### Reverse Proxy

Define first _downstream_s, which are the systems 
behind the reverse proxy to which you are providing access.

```yaml
downstream:
  - target: "system-name" # a name recognisable by you for the target system
    protocol: "http" # used by the target system
    host: "localhost" # at which the target system can be reached
    port: 8080 # that the target system exposes
    base: "/" # the base part of the URL the target system expects
    path-mapper: 
    # how to map a path on which a request is received to the request made to the target system
      type: forward-unchanged # explained later
```

The path at which `ferp` should receive requests which should
be proxied to the downstream system is not defined here - 
the downstream section only contains data about the target systems.

The _http_ and _https_  _incoming_ sections define from which paths
`ferp` should receive requests 
and forward them to which downstream systems.

```yaml
http:
  port: 80 # where should the HTTP reverse proxy server be exposed
  incoming:
    - path: "/static/index.html" # on a request to /static/index.html to the reverse proxy
      methods:
        - "GET" # if it's a GET method
        - "POST" # or a post method
      target: "system-name" # forward to the downstream with this target value
https:
  port: 443 # where should the HTTPS reverse proxy server be exposed
  cert-file: "/path/to/fullchain.pem" # where to find this server's cert file
  key-file: "/path/to/privkey.pem" # where to find this server's private key file
  incoming:
    - path: "/api/*" # on a request to anything starting with /api/ to the reverse proxy
      methods:
        - "ALL" # for any method
      target: "system-name" # forward to the downstream with this target value
```

#### Path Mapping

The link between the path on which the reverse proxy receives
a request and the path in the URL to which it forwards the
request to the downstream system is made in the path mapper.

Consider the above example downstream, with a base of `/`,
and the http incoming, on path `/static/index.html`.

The path mapper configured for the downstream takes the
incoming path (`/static/index.html`) transforms it _somehow_,
appends the result to the base `/` (taking care not to 
end up with e.g. `//static/...`), and uses the resulting
path in the forwarded request.

The transformation depends upon the chosen path mapper.
There are two "built in" path mappers.

##### Forward Unchanged

```yaml
    path-mapper:
      type: "forward-unchanged"
```

The reverse proxy simply re-uses the path at which it 
received the request. Hence `/static/index.html` becomes
`/static/index.html`, is appended to `/`, one slash 
is dropped, and the request is made to the path
`/static/index.html` on the downstream system.

##### Remove Prefix

```yaml
    path-mapper:
      type: "remove-prefix"
      prefix: "/static/"
```

The reverse proxy attempts to remove the configured prefix
from the start of the path on which it received the request.
Hence `/static/index.html` becomes `index.html`, 
is appended to `/`, and the request is made 
to the path `/index.html` on the downstream system.

Note: if the path is not prefixed with the configured prefix,
the path will not be changed and this mapper will behave
the same way as forward-unchanged.

### Redirects

It is often useful/convenient to define aliases/shortened urls
that point to certain parts of one's application. This shouldn't
really be a concern of the service running behind the proxy -
it shouldn't really know about the "outside world" - so they
are more naturally handled by the reverse proxy itself.

These may be configured on the HTTP or HTTPS server, and may
only redirect to the same server (i.e. the HTTP server cannot
be configured to redirect to the HTTPS server, nor vice versa).

```yaml
http:
  # ...
  redirects:
    - from: "/home" # on a request to /home on the http server
      to: "/static/index.html" # redirect to this path on the http server
      methods:
        - "GET" # if the method is GET
```

### Ordering

Routes follow ordering/preference rules you would expect
for any HTTP(S) server software or code. Routes are configured
in the order in which they appear in the configuration (first to 
last -> top to bottom) and are matched in this order.
Hence e.g. in this configuration

```yaml
https:
  # ...
  incoming:
    - path: "/api/*"
      methods:
        - "ALL"
      target: "api-server"
    - path: "/api/special/route"
      methods:
        - "GET"
      target: "special-api-server"
```

the `special-api-server` target can never be reached, because
`/api/special/route` will always be matched by the previous
_incoming_ with path `/api/*`.
