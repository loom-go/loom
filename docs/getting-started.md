---
title: "GETTING STARTED"
weight: 2
---

We'll cover: how to install loom, how to setup a renderer, and how to run your first loom app!

---

### Installing Loom

In a new folder, initialize a new project:

```bash {style=tokyonight-moon}
go mod init my-project
```

And install loom:

```bash {style=tokyonight-moon}
go get github.com/loom-go/loom
```

<br/>

### Installing a renderer

When setting up a new loom app, you must choose a [renderer](/docs/concepts#renderer). This is where you decide which plateform your loom app will run on.

[LOOM-TERM](/term), for building Terminal UIs.
<br/>
[LOOM-WEB](/web), for building Web SPAs.

How to intall each:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```bash {style=tokyonight-moon}
go get github.com/loom-go/term
```

> LOOM-TERM uses CGO. Make sure you have a C compiler like GCC or MinGW installed.

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
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for t := range ticker.C {
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
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for t := range ticker.C {
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

    "github.com/loom-go/term"
)

func main() {
	app := term.NewApp()

	for err := range app.Run(RenderFullscreen, App) {
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
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for t := range ticker.C {
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

	for err := range app.Run(RenderFullscreen, App) {
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
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for t := range ticker.C {
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

- create index.html file
- wasm_exec.js
- build
- open index.html

{{< /tab >}}
{{< /tabs >}}

**GREAT SUCCESS**

---

From there it's up to you!

Be sure to have a look at -> [CORE CONCEPTS](/docs/concepts) to understand more about loom.

If you’re coming from a signal-based JavaScript framework, make sure you have a quick read of -> [SIGNALS SCHEDULING](/docs/guides/reactivity#scheduling) and -> [BINDING](/docs/guides/binding) to understand the differences with loom.
