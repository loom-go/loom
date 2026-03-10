---
title: "EXAMPLES"
weight: 10
---

#### Counter

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
func Counter() Node {
	count, setCount := Signal(0)

    go func(self Component) {
        for !self.IsDisposed() {
            time.Sleep(time.Second / 30)
            setCount(count() + 1)
        }
    }(Self())

	return P(Text("Count: "), BindText(count))
}
```

<video src="/medias/counter-term.mp4" autoplay loop muted playsinline></video>

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
func Counter() Node {
	count, setCount := Signal(0)

    go func(self Component) {
        for !self.IsDisposed() {
            time.Sleep(time.Second / 30)
            setCount(count() + 1)
        }
    }(Self())

	return P(Text("Count: "), BindText(count))
}
```

<video src="/medias/counter-web.mp4" autoplay loop muted playsinline></video>

{{< /tab >}}
{{< /tabs >}}

<br/>

#### Conditions

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
func Condition() Node {
    display, setDisplay := Signal(false)

    toggle := func(*term.EventMouse) {
        setDisplay(!display())
    }

    return Box(
        Box(Text("toggle"), Apply(On{Click: toggle})),

        Show(display, func() Node {
            return Text("am i visible now?")
        }),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
func Condition() Node {
    display, setDisplay := Signal(false)

    toggle := func(*web.EventMouse) {
        setDisplay(!display())
    }

    return Div(
        Box(Text("toggle"), Apply(On{Click: toggle})),

        Show(display, func() Node {
            return Text("am i visible now?")
        }),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

<br/>

#### Lists

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
func FruitList() Node {
    fruits, setFruits := Signal([]string{"banana", "apple", "orange"})

    return Box(
        For(fruits, func(fruit Accessor[string], index Accessor[int]) Node {
            return P(BindText(fruit))
        }),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
func FruitList() Node {
    fruits, setFruits := Signal([]string{"banana", "apple", "orange"})

    return Ul(
        For(fruits, func(fruit Accessor[string], index Accessor[int]) Node {
            return Li(BindText(fruit))
        }),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

<br/>
