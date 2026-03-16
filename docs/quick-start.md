---
title: "QUICK START"
weight: 2
---

We'll cover: how to install loom, how to setup a renderer, and how to run your first loom app!

---

### Setup

In a new folder, initialize a new project:

```bash {style=tokyonight-moon}
go mod init my-project
```

Install loom:

```bash {style=tokyonight-moon}
go get github.com/loom-go/loom
```

And install a [renderer](/docs/get-started/concepts#renderer):

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

<br/>

### Creating a component

In a new `main.go` file:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
package main

import (
	"time"

    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    . "github.com/loom-go/term/components"
)

var (
	styleContainer = Style{
		Width:          "100%",
		Height:         "100%",
		AlignItems:     "center",
		JustifyContent: "center",
	}

	styleTime = Style{
		Color: "#6ac482",
	}
)

func App() loom.Node {
    now, setNow := Signal(time.Now())

    go func(self loom.Component) {
        for t := range time.Tick(time.Second) {
            if self.IsDisposed() {
                return
            }

            setNow(t)
        }
    }(Self())

    return Box(
        BindText(now, Apply(styleTime)),
        Apply(styleContainer),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
package main

import (
	"time"

    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    . "github.com/loom-go/web/components"
)

var (
	styleContainer = Style{
		"width":           "100vw",
		"height":          "100vh",
        "display":         "flex",
		"align-items":     "center",
		"justify-content": "center",
	}

	styleTime = Style{
		"color": "#6ac482",
	}
)

func App() loom.Node {
    now, setNow := Signal(time.Now())

    go func(self loom.Component) {
        for t := range time.Tick(time.Second) {
            if self.IsDisposed() {
                return
            }

            setNow(t)
        }
    }(Self())

    return Div(
        BindText(now, Apply(styleTime)),
        Apply(styleContainer),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

<br/>

### Using the renderer

In the same `main.go` file:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    // ...

    "fmt"
    "os"
    "github.com/loom-go/term"
)

func main() {
	app := term.NewApp()

	for err := range app.Run(term.RenderFullscreen, App) {
		app.Close()

		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    // ...

    "github.com/loom-go/web"
)

func main() {
	app := web.NewApp()

	for err := range app.Run("#root", App) {
        web.ConsoleError(fmt.Sprintf("Error: %v\n", err)
	}
}
```

{{< /tab >}}
{{< /tabs >}}

<br/>

### Running the app

<br/>
<details>
<summary>Full code</summary>

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
package main

import (
    "fmt"
    "os"
	"time"

    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    "github.com/loom-go/term"
    . "github.com/loom-go/term/components"

)

var (
	styleContainer = Style{
		Width:          "100%",
		Height:         "100%",
		AlignItems:     "center",
		JustifyContent: "center",
	}

	styleTime = Style{
		Color: "#6ac482",
	}
)

func App() loom.Node {
    now, setNow := Signal(time.Now())

    go func(self loom.Component) {
        for t := range time.Tick(time.Second) {
            if self.IsDisposed() {
                return
            }

            setNow(t)
        }
    }(Self())

    return Box(
        BindText(now, Apply(styleTime)),
        Apply(styleContainer),
    )
}

func main() {
	app := term.NewApp()

	for err := range app.Run(term.RenderFullscreen, App) {
		app.Close()

		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
package main

import (
	"time"

    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    "github.com/loom-go/web"
    . "github.com/loom-go/loom-web/components"
)

var (
	styleContainer = Style{
		"width":           "100vw",
		"height":          "100vh",
        "display":         "flex",
		"align-items":     "center",
		"justify-content": "center",
	}

	styleTime = Style{
		"color": "#6ac482",
	}
)

func App() loom.Node {
    now, setNow := Signal(time.Now())

    go func(self loom.Component) {
        for t := range time.Tick(time.Second) {
            if self.IsDisposed() {
                return
            }

            setNow(t)
        }
    }(Self())

    return Div(
        BindText(now, Apply(styleTime)),
        Apply(styleContainer),
    )
}

func main() {
	app := web.NewApp()

	for err := range app.Run("#root", App) {
        web.ConsoleError(fmt.Sprintf("Error: %v\n", err)
	}
}
```

{{< /tab >}}
{{< /tabs >}}

</details>

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```bash {style=tokyonight-moon}
go run main.go
```

And you should see the current time in fullscreen!

{{< /tab >}}
{{< tab >}}

Create an `index.html` file:

```html {style=tokyonight-moon}
<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>App</title>
  </head>

  <body>
    <div id="app"></div>

    <script src="https://cdn.jsdelivr.net/gh/golang/go@go1.25.0/lib/wasm/wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(
        fetch("main.wasm"),
        go.importObject,
      ).then((r) => go.run(r.instance));
    </script>
  </body>
</html>
```

Build your `main.go` in wasm:

```go {style=tokyonight-moon}
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

And serve your files with your favorite http server!
Here with [`serve`](https://www.npmjs.com/package/serve):

```go {style=tokyonight-moon}
npx serve
```

{{< /tab >}}
{{< /tabs >}}

**GREAT SUCCESS**

---

From there it's up to you!

Be sure to have a look at -> [CORE CONCEPTS](/docs/get-started/concepts) to understand more about loom.

If you’re coming from a signal-based JavaScript framework, make sure you have a quick read of -> [SIGNALS SCHEDULING](/docs/guides/reactivity#scheduling) and -> [BINDING](/docs/guides/binding) to understand the differences with loom.
