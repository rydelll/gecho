<div align="center">
    <a href="https://github.com/rydelll/gecho">
        <picture>
            <img alt="gecko" src="docs/image/gecko.png" height="128">
        </picture>
    </a>
    <h1>Gecho</h1>

<a href="https://go.dev/doc/devel/release"><img alt="go" src="https://img.shields.io/github/go-mod/go-version/rydelll/gecho"></a>
<a href="https://github.com/rydelll/gecho/actions"><img alt="workflow" src="https://github.com/rydelll/gecho/actions/workflows/ci.yml/badge.svg"></a>
<a href="https://github.com/rydelll/gecho/blob/main/LICENSE"><img alt="license" src="https://img.shields.io/github/license/rydelll/gecho"></a>
<a href="https://github.com/rydelll/gecho/issues"><img alt="issues" src="https://img.shields.io/github/issues/rydelll/gecho.svg"></a>
</div>

Gecho is a TCP echo serve written in Go that listens for connection requests, reads the data it recieves from the client, and writes it back unmodified.

# Install

To install and start Gecho, run the following commands.

```bash
go install github.com/rydelll/gecho@latest
gecho -port 7777
```
