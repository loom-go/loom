---
title: "SETUP"
weight: 1
---

#### Tooling

Make sure you have a recent [Go](https://go.dev/doc/install) version installed (>= v1.23).

#### 1. Installing loom

Loom works like any other Go module. It can be installed with `go get` in any project:

```bash {style=tokyonight-moon}
go get github.com/loom-go/loom
```

#### 2. Installing a renderer

When setting up a new loom app, you must choose a [renderer](/docs/get-started/concepts#renderer). This is where you decide which plateform your loom app will run on.

[\*] <a href="/term/intro">LOOM-TERM -></a> | For building Terminal UIs.<br/>
[\*] <a href="/web/intro">LOOM-WEB -></a> | For building Web SPAs.

How to intall each:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```bash {style=tokyonight-moon}
go get github.com/loom-go/term
```

{{< callout type="warning" >}}
**LOOM-TERM** uses [CGO](https://go.dev/blog/cgo). Make sure you have a C compiler installed like GCC on Linux/Darwin, or MinGW on Windows.
{{< /callout >}}

{{< /tab >}}
{{< tab >}}

```bash {style=tokyonight-moon}
go get github.com/loom-go/web
```

{{< /tab >}}
{{< /tabs >}}

---

Keep reading -> [CORE CONCEPTS](/docs/get-started/concepts) to understand more about loom and how it works.
