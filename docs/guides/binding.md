---
title: "BINDING"
weight: 2
---

Everything you need to know about _binding_ in loom.

This document expects you to be familiar with loom's core concepts and signal-based reactivity. If not -> [CORE CONCEPTS](/docs/concepts), and -> [REACTIVITY](/docs/guides/reactivity)

---

### What's binding?

Binding is the act of assigning a reactive signal to the UI.

It syncs one or more parts of the UI from the signals it depends on. When a signal changes, the UI updates.

For instance, _binding_ a text element's content from a signal's value means: whenever the signal changes, the content of the element gets updated with the value of the signal. It syncs the value of the element, with the value of the signal.

### Binding in loom

If you're coming from a JavaScript framework with JSX, you might be used to binding being implicit. When you use a signal in your JSX it binds that signal automatically to that element. This is not the case in loom.

In loom, binding is explicit. It's on you to decide what part of the tree updates or not.

For loom to be pure Go without any compilation overhead, this is a trade off that needed to be made. But I promise you'll get used to it, and you might even end up preferring it!

```go {style=tokyonight-moon}
count, setCount := Signal(0)

return P(
    // manually binding Text() to the `count` signal
    // (we'll learn how to simplify this, it's just a demo)
    Bind(func() Node {
        str := fmt.Sprintf("Count: %d", count())
        return Text(str)
    }),
)
```

Explicit binding gives you more control over the tree and how it reacts to changes. You can update precisely what's needed. From a single attribute, to a bigger part of the UI, it gives you full control without unnecessary costs.

#### Bind()

As shown above, [`Bind()`](/docs/components/bind) is the default way to _bind_ a signal to the tree. It takes in a function returning a node that be called each time the signal changes.

```go {style=tokyonight-moon}
fruits, setFruits := Signal([]string{"banana", "apple"})

return Bind(func() Node {
    // reading the `fruits` signal.
    // making it a dependency of this Bind.
    length := len(fruits())

    // each time `fruits` changes, this function will be called to recompute the tree.
    // here we're just updating the text.
    // so the Text element will simply update its content.

    if length == 0 {
        return Text("Zero fruit")
    }
    if length == 1 {
        return Text("1 fruit")
    }

    return Text(fmt.Sprintf("%d fruits", length))
    // could also be written:
    // return Fragment(Text(length), Text(" fruits"))
})
```

From there it's up to you! You can use `Bind()` for micro-updating specific parts of the tree like shown above, or for reconstructing a whole chunk of your UI.

#### BindX()

Most components comes with a `Bind()` wrapper to make it easier for you to update it.

`Text()` has `BindText()`. `Attr()` has `BindAttr()`, `Style{}` has `BindStyle{}` etc.

They essensiatlly are just wrappers around the standard component, but takes a signal (a function returning a value) instead of the value directly:

{{< tabs items="BindText, Bind" >}}
{{< tab >}}

```go {style=tokyonight-moon}
content, setContent := Signal("")

return P(BindText(content))
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
content, setContent := Signal("")

return P(
    Bind(func() Node {
        return Text(content())
    }),
)
```

{{< /tab >}}
{{< /tabs >}}

{{< tabs items="BindAttr, Bind" >}}
{{< tab >}}

```go {style=tokyonight-moon}
value, setValue := Signal("")

return Input(BindAttr("value", value))
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
value, setValue := Signal("")

return Input(
    Bind(func() Node {
        return Attr("value", value())
    }),
)
```

{{< /tab >}}
{{< /tabs >}}

{{< tabs items="BindStyle, Bind" >}}
{{< tab >}}

```go {style=tokyonight-moon}
color, setColor := Signal("#777")

return Box(
    Text("a box"),

    Apply(
        Style{Width: 100, Height: 100},
        BindStyle{BackgroundColor: color}, // taking signals instead of values
    ),
)
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
color, setColor := Signal("#777")

return Box(
    Text("a box"),

    Apply(Style{Width: 100, Height: 100}),
    Bind(func() Node {
        return Apply(Style{
            BackgroundColor: color()
        })
    }),
)
```

{{< /tab >}}
{{< /tabs >}}
