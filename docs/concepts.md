---
title: "CORE CONCEPTS"
weight: 3
---

Loom has five core concepts.

They are the roots of what makes loom a composable and reusable framework.

---

### Node

Nodes are the most fundamental primitives. They represent each and every parts of your UI.

A Node is responsible for displaying, updating, and destroying one or more piece of interface. They are the synchronisation between your code, and what's displayed on the screen.

But you will most likely never write a Node yourself. They are low-level and mostly spesific to [renderers](#renderer). Instead you will interact with higher-level abtractions built on top of them ([components](#component)).

---

### Component

Components are an abstraction on top of [nodes](#node).

It is nothing but a regular Go function, that returns a Node:

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
)

func MyComponent() loom.Node {
    return nil
}
```

Loom (and the [renderer](#renderer) of your choice) provide various built-in components. Theses components can be used and composed in your own components to build a complete UI.

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    . "github.com/loom-go/term/components"
)

func MyComponent() loom.Node {
    return Fragment( // using the Fragment() component from loom
        P(Text("hello")), // P() and Text() from the loom-term renderer

        P(
            Text("hello in pink"),
            Apply(Style{ // also from loom-term
                BackgroundColor: "#ffc0cb"
            }),
        ),
    ),
}
```

Components are only called at mount. Meaning you can spin up goroutines for async work and update states from there!

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    . "github.com/loom-go/term/components"
)

func Counter() Node.loom. {
    // if your not sure what a signal is,
    // read the following section about reactivity
	count, setCount := Signal(0)

    go func() {
        for {
            time.Sleep(time.Second)
            setCount(count() + 1)
        }
    }()

	return P(Text("Count: "), BindText(count))
}
```

---

### Applier

Appliers are similar to attributes in HTML. They _apply_ something on a Node.

They come as Go struct that you instantiate yourself,
and apply on a Node with loom's [`Apply()`](docs/components/apply) component.

For instance the [`Style{}`](/term/appliers/style) applier from [LOOM-TERM ->](/term) :

```go {style=tokyonight-moon}
// instantiate a Style{} applier
var styleBox = Style{
    Width: 10,
    Height: 10,
    BackgroundColor: "red",
}

return Box(
    Text("a box"),

    Apply(styleBox), // apply styleBox on the Box()
)
```

Renderers can also provide extended `Apply()` components like LOOM-TERM's [`ApplyOn()`](/term/components/applyon):

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    "github.com/loom-go/term"
    . "github.com/loom-go/term/components"
)

var (
	styleBox        = Style{BackgroundColor: "red"}
	styleBoxHover   = Style{BackgroundOpacity: 0.5}
)

func MyComponent() loom.Node {
	return Box(
        Apply(styleBox),
        ApplyOn("hover", styleBoxHover),
	)
}
```

<details>
<summary>
Example with
<a href="/docs/appliers/ref"><code>Ref{}</code></a>,
<a href="/term/appliers/on"><code>On{}</code></a>,
<a href="/term/appliers/style"><code>Style{}</code></a> and
<a href="/term/components/applyon"><code>ApplyOn()</code></a>
</summary>

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    "github.com/loom-go/term"
    . "github.com/loom-go/term/components"
)

var (
	styleInput    = Style{Width: 30, BackgroundColor: "lightgray"}
	styleBtn      = Style{BackgroundColor: "gray"}
    styleBtnHover = Style{BackgroundOpacity: 0.5}
)

func MyComponent() loom.Node {
	var input term.InputElement

	focus := func(*term.EventMouse) { input.Focus() }
	blur := func(*term.EventMouse) { input.Blur() }

	return Box(
		// apply can take multiple appliers
		Input(Apply(
            Ref{Ptr: &input},
            styleInput,
        )),

		Box(
            Text("focus"),
            Apply(On{Click: focus}, styleBtn),
            ApplyOn(styleBoxHover),
        ),
		Box(
            Text("blur"),
            Apply(On{Click: blur}, styleBtn),
            ApplyOn(styleBoxHover),
        ),
	)
}
```

</details>

---

### Reactivity

Reactivity makes the UI respond to changes.

It can be updating a color when a user clicks a button, or refreshing a list when a user fills an input, or anything else related to a reaction from change.

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
    . "github.com/loom-go/loom/components"
    . "github.com/loom-go/term/components"
)

func MyComponent() loom.Node {
    text, setText := Signal("")

    update := func(e *EventInput) {
        setColor(e.InputValue()) // update the text with what the user typed
    }

    return Fragment(
        // note the use of BindText().
        // reactivity is explicit in loom.
        // read the BINDING guide to learn more
        P(Text("You typed: "), BindText(text)),

        InputText(On("input", udpate)),
    ),
}
```

As shown above, in loom reactivity is signal-based. If you're coming from a modern JavaScript framework, you'll feel right at home.

But that doesn't mean reactivity works exactly the same as in JS frameworks. Make sure to read -> [SIGNALS SCHEDULING](/docs/guides/reactivity#scheduling) and -> [BINDING](/docs/guides/binding) to understand the differences with loom.

Or if you want to understand more about using reactivity in general, you can read the full guide -> [REACTIVITY](/docs/guides/reactivity)

---

### Renderer

By itself, loom cannot display anything on your screen. It needs a Renderer for that.

A Renderer is responsible for displaying content on screen by providing you plateform-specific components to build a UI with.<br/>
For instance a web renderer would provide DOM components like `<div>` or `<ul>`. And a theoretical mobile renderer would provide native components like `View`, `Text`, `ScrollView`, etc.

**There's currently two official renderers:**<br/>
[\*] <a href="/term/intro">LOOM-TERM -></a> | For building Terminal UIs.<br/>
[\*] <a href="/web/intro">LOOM-WEB -></a> | For building Web SPAs.

If you'd like to get started with one of the two -> [GET STARTED](/docs/getting-started)
