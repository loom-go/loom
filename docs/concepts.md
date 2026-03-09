---
title: "CORE CONCEPTS"
weight: 3
---

Loom has four core concepts.

They are the roots of the architecture that makes loom a composable and reusable framework that can be used for any plateform.

---

### Node

Nodes are the most fundamental primitives. They represent each and every parts of your UI.

A Node is responsible for displaying, updating, and destroying one or more parts of your UI. They are the synchronisation between your code, and what's displayed on the screen.

But you will most likely never write a Node yourself. They are low-level and mostly spesific to [renderers](#renderer). Instead you will interact with higher-level abtractions built on top of them ([components](#component)).

---

### Component

Components are an abstraction on top of [nodes](#node).

It is nothing but a regular Go function, that returns a Node:

```go {style=tokyonight-moon}
import (
    . "github.com/AnatoleLucet/loom"
)

func MyComponent() Node {
    return nil
}
```

Loom (and the [renderer](#renderer) of your choice) provide various built-in components. Theses components can be used and composed in your own components to build a complete UI.

```go {style=tokyonight-moon}
import (
    . "github.com/AnatoleLucet/loom"
    . "github.com/AnatoleLucet/loom/components"
    . "github.com/AnatoleLucet/loom-term/components"
)

func MyComponent() Node {
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
    . "github.com/AnatoleLucet/loom"
    . "github.com/AnatoleLucet/loom/components"
    . "github.com/AnatoleLucet/loom-term/components"
)

func Counter() Node {
    // if your not sure what a signal is,
    // read the next section about reactivity
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

### Reactivity

It is what makes your UI react to changes.

It can be updating a color when a user clicks a button, or refreshing a list when a user fills an input, or anything else related to a reaction from change.

```go {style=tokyonight-moon}
import (
    . "github.com/AnatoleLucet/loom"
    . "github.com/AnatoleLucet/loom/components"
    . "github.com/AnatoleLucet/loom-term/components"
)

func MyComponent() Node {
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

If you're coming from a JS framework, make sure you read -> [SIGNALS SCHEDULING](/docs/guides/reactivity#scheduling) and -> [BINDING](/docs/guides/binding) to understand the differences with loom.

Or if you want to understand more about using reactivity, you can read the full guide -> [REACTIVITY](/docs/guides/reactivity)

---

### Renderer

By itself, loom cannot display anything on your screen. It needs a Renderer for that.

A Renderer is responsible for displaying content on screen by providing plateform-specific components for the use to build a UI with. For instance a web renderer would provide DOM components like \<div\> or \<ul\> to the user. While a theoretical mobile renderer would provide native components for View, Text, ScrollView, etc.

**There's currently two official renderers:**

[*] [LOOM-TERM ->](/term/intro) | For building Terminal UIs.

[*] [LOOM-WEB ->](/web/intro) | For building Web SPAs.

<br/>

If you'd like to get started with one of the two -> [GET STARTED](/docs/getting-started)
