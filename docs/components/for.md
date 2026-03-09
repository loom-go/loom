---
weight: 3
title: For()
---

```go {style=tokyonight-moon}
func For[T any](
    items func() []T,
    keyer func(item T) any,
    mapper func(item Accessor[T], index Accesor[T]) Node,
) Node
```
