---
title: "INTRO"
weight: 1
---

Loom is a framework for building user interfaces, for any plateform.

You define declarative components written in pure Go, and loom renders them using a plateform-specific [renderer](/docs/concepts.md#renderer).

{{< tabs items="TERM,WEB" >}}

{{< tab >}}

```go {style=tokyonight-moon}
// define your component
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

func main() {
	app := term.NewApp()

    // render it using our terminal renderer
    errs := app.Run(term.RenderInline, Counter)
}
```

<video src="/medias/counter-term.mp4" autoplay loop muted playsinline></video>

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
// define your component
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

func main() {
	app := web.NewApp()

    // render it using our web renderer
    errs := app.Run("#root", Counter)
}
```

<video src="/medias/counter-web.mp4" autoplay loop muted playsinline></video>

{{< /tab >}}
{{< /tabs >}}

If you're coming from JavaScript, this should feel very familiar.

You can use it to build and compose as many components as your UI needs. From simple static components, to complex UI layouts with hundreds or thousands of moving parts.

---

If you'd like to get started -> [GETTING STARTED](/docs/getting-started)

Or to keep reading about loom and how it works -> [CORE CONCEPTS](/docs/concepts)
