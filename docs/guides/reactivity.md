---
title: "REACTIVITY"
weight: 1
---

An introduction to reactivity and signals.

This document expects you to be familiar with loom's core concepts. If not -> [CORE CONCEPTS](/docs/concepts)

---

> If you're coming from a signal-based JavaScript framework like SolidJS or Svelte, you will rapidly notice some [divergence](#scheduling) in the reactive model.
> But you'll get used to it very soon, I promise.

### Signal

Signals are the core primitive to reactivity.

Like a variable, a signal holds a value. This value can be read, updated or kept as is.
But unlike a variable, signals can be subscribed to. When its value changes, the subscribers get notified.

Let's look at an example:

```go {style=tokyonight-moon}
// define a signal called "count" with a default value of 0
count, setCount := Signal(0)

// `count` can read the value
count() // 0

// `setCount` can update it
setCount(1)
count() // 1


// example: increment by 1
setCount(count() + 1)
```

From here, there's no reason to use a signal over a regular variable. To see where signals become usefull, we must dive into subscribers.

### Memo

A Memo is one of multiple type of subscriber. It listens to one or more signals, and gets notified when there's a change.

It is responsible for holding a value, and recomputing that value when one of its dependencies (the signals it depends one) changes.

```go {style=tokyonight-moon}
// define the same signal as before
count, setCount := Signal(0)

// but this time declare a memo
double := Memo(func() int {
    // becaues we're readding `count` here,
    // this function will be called each time `count` changes.
    // this is the subscriber, and `count` is the dependency

    // read `count` and double its value
    return count() * 2
})

count() // 0
double() // 0

// update count (notifying the memo)
setCount(2)

count() // 2
double() // 4, the memo's value got updated!
```

### Effect

An effect is another type of subscriber.

In loom it behaves exactly like a Memo. It takes a function, and runs that function synchronously when a dependency changes.
Except its purpose is not to recompute a value, but to synchronize an external system.

```go {style=tokyonight-moon}
count, setCount := Signal(0)

Effect(func() {
    // printing to the terminal each time `count` changes
    fmt.Println("count:", count())
})

setCount(10) // prints "count: 10" to the terminal
setCount(22) // prints "count: 22" to the terminal
```

### Binding

Binding is the act of assigning a reactive signal to a part of the UI. For instance, _binding_ a text element's content to a signal's value means: whenever the signal changes, the text element's content gets updated with that value.

If you're coming from a JavaScript framework, you might be used to implicit binding in JSX. But this is not the case with loom.

Since loom is pure Go and without any compilation overhead, binding is explicit. _You_ decide what part of the UI updates.

```go {style=tokyonight-moon}
func Counter() Node {
	count, setCount := Signal(0)

    go func() {
        for {
            time.Sleep(time.Second)
            setCount(count() + 1)
        }
    }()

	return P(
        Text("Count: "),
        // BindText() takes a signal and returns a text Node.
        // this Node gets updated each time `count` changes.
        BindText(count),
    )
}
```

At first it might seem like a downgrade from JSX, but with time you will most likely see the benefits. To read more about binding -> [BINDING](/docs/guides/binding)

From there you've covered 3/4 of what reactivity is in loom! The rest are conveniences for this paradigm, and best practicies.

If you want to understand more you can read the following section [about scheduling](#scheduling), dive into [the references](http://localhost:1313/docs/reactivity/signal/), or take a look at [the reactive model](https://github.com/AnatoleLucet/sig) that was built to power loom.

---

### Scheduling

> If you're coming from a signal-based JavaScript framework like SolidJS or Svelte 5, this is where loom's reactive model diverges a bit.

Due to how Go's internal scheduling works, and to make signals work across goroutines boundaries, the reactive system is fully synchronous.

Meaning two things:

\- When you update a signal, it synchronously calls every subscribers (memos, effects, etc.) inside that `setCount()` call.
No automatic batching of signal updates like in Svelte 5 or SolidJS 2, everything is instantaneous and synchronous.

\- Effects initialize synchronously. An effect initialization is not deferred, it runs at the moment you call `Effect(fn)`.
If you're coming from JavaScript, effects in loom behave like SolidJS's `createRenderEffect()` or Svelte 5's `$effect.pre()`, they are blocking and run instantly.
But [that doesn't mean there's no prioritization in loom](https://github.com/AnatoleLucet/sig?tab=readme-ov-file#features).

An example will speak for itself:

```go {style=tokyonight-moon}
// just a regular signal
count, setCount := Signal(0)

// double depends on count()
double := Memo(func() int {
    fmt.Println("running double:", count())
    return count() * 2
})

// and effect depends on double()
Effect(func() {
    fmt.Println("running effect:", double())
})

fmt.Println("initialized")

// the above has the following dependency chain: `count <- double <- effect`.
// the effect depends on double, and double depends on count.

// for now, we've only initialized our signal and subscribers.
// the terminal prints:
// "running double: 0"
// "running effect: 0"
// "initialized" <- notice how this runs after, because everything is synchronous


// now lets udpate our signal:
fmt.Println("preparing to update")
setCount(10)
fmt.Println("updated successfully")

// the terminal now prints:
// "preparing to udpate"
// "running double: 10"
// "running effect: 20"
// "updated successfully" <- notice how the effect and memo ran before this. they executed inside the `setCount()` call

// meaning if you call setCount twice, the memo and effect will run twice:
setCount(11)
setCount(12)

// "running double: 11"
// "running effect: 22"
// "running double: 12"
// "running effect: 24"

// this is where Batch(fn) comes into play!
// because batching is manual in loom, not automatic.
```

If this is your first encounter with a reactive system, it might seem like a very obvious behavior.
But it is certainly not for users coming from an ~asynchronous reactive system like SolidJS's or Svelte's, or event ReactJS's.
