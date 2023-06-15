# proxy

This is a really simple Go proxy to let you access another website
through a proxy. The usage is really simple:

```
usage: proxy <url> [<address>]
Proxy <url> onto <address>
    url            the url to proxy to
    address        the address to serve to, default: localhost:3000

example:
    proxy https://github.com 0.0.0.0:80
```
