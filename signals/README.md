<h1 align="center"><code>signals</code></h1>

<p align="center">Loom's signal-based reactivity model, powered by <a href="https://github.com/AnatoleLucet/sig"><code>sig</code></a>.</p>

```go
count, setCount := Signal(0)

Effect(func() {
    fmt.Println("changed", count())
})

setCount(10)
```

## Usage

### Installation

```bash
go get github.com/AnatoleLucet/loom/signals
```

### Basic counter

```go
package main

import (
    "fmt"

    . "github.com/AnatoleLucet/loom/signals"
)

func main() {
    count, setCount := Signal(0)

    Effect(func() {
        fmt.Println("changed:", count())
    })

    setCount(count() + 1)
}
```

### API

#### `Signal`

```go
count, setCount := Signal(0)
fmt.Println(count()) // 0

setCount(10)
fmt.Println(count()) // 10
```

#### `Memo`

```go
count, setCount := Signal(1)
double := Memo(func() int {
    return count() * 2
})
fmt.Println(double()) // 2

setCount(10)
fmt.Println(double()) // 20
```

#### `Effect`

```go
count, setCount := Signal(1)

Effect(func() {
    fmt.Println("changed:", count())
})

setCount(10) // changed: 10
```

#### `Batch`

```go
count1, setCount1 := Signal(1)
count2, setCount2 := Signal(2)

Effect(func() {
    fmt.Println("changed:", count1(), count2())
})

Batch(func() {
    setCount1(10)
    setCount2(11)
}) // changed: 10, 11
```

#### `Untrack`

```go
count, setCount := Signal(1)

Effect(func() {
    fmt.Println("changed:", Untrack(count))
})

setCount(10) // *nothing*
```
