<h1 align="center">「#」</h1>

<p align="center">A reactive component framework for TUIs, the Web, and more.</p>

```go
func Counter() Node {
	count, setCount := Signal(0)

	go func() {
		for {
			time.Sleep(time.Second / 30)
			setCount(count() + 1)
		}
	}()

	return P(Text("Count: "), BindText(count))
}
```

## Features

- **Pure Go** | No extra compiler.
- **Multi-plateform** | Built-in support for TUIs and SPAs.
- **Signal-based** | Concurrency-safe [reactive model](https://github.com/loom-go/loom/tree/main/signals/README.md) with signals, effects, memos, etc.
- **Components** | Define your UI as declarative JSX-like components.

## Quick-start

```bash
go mod init my-project
go get github.com/loom-go/loom github.com/loom-go/term
```

```go
package main

import (
	"log"
	"time"

	"github.com/loom-go/loom"
	. "github.com/loom-go/loom/components"
	"github.com/loom-go/term"
	. "github.com/loom-go/term/components"
)

func Counter() loom.Node {
	frame, setFrame := Signal(0)

	go func() {
		for {
			time.Sleep(time.Second / 120)
			setFrame(frame() + 1)
		}
	}()

	return Box(Text("Count: "), BindText(frame))
}

func main() {
	app := term.NewApp()

	for err := range app.Run(term.RenderInline, Counter) {
		app.Close()
		log.Fatalf("Error: %v\n", err)
	}
}
```

```bash
go run .
```

And it's live!

## Documentation

You can visite [loom's website](https://loomui.dev) for the full documentation.

## License

[MIT](./LICENSE)
