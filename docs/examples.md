---
title: "EXAMPLES"
weight: 10
---

Fully functioning examples can be found at [`github.com/go-loom/loom/examples`](https://github.com/loom-go/loom/tree/main/examples).

#### Counter

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
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

<video src="/medias/counter-term.mp4" autoplay loop muted playsinline></video>

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
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

<video src="/medias/counter-web.mp4" autoplay loop muted playsinline></video>

{{< /tab >}}
{{< /tabs >}}

---

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

---

#### Lists

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
func FruitList() Node {
    fruits, setFruits := Signal([]string{"banana", "apple", "orange"})

    return Box(
        For(fruits, func(fruit string, index Accessor[int]) Node {
            return P(Text(fruit))
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
        For(fruits, func(fruit string, index Accessor[int]) Node {
            return P(Text(fruit))
        }),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

---

#### Goroutine cancellation

See [`Self()`](/docs/reactivity/self).

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
func MyComponent() Node {
    go func(self Component) {
        for {
            select {
            case <-self.Disposed():
                // stop Goroutine when component is diposed
                return
            default:
            }

            // keep looping
        }
    }(Self())

	return Text("My component")
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
func MyComponent() Node {
    go func(self Component) {
        for {
            select {
            case <-self.Disposed():
                // stop Goroutine when component is diposed
                return
            default:
            }

            // keep looping
        }
    }(Self())

	return Text("My component")
}
```

{{< /tab >}}
{{< /tabs >}}

---
