---
weight: 5
title: Keyed()
---

```go {style=tokyonight-moon}
func Keyed[T any](
    items func() []T,
    keyer func(item T) any,
    mapper func(item Accessor[T], index Accesor[int]) Node,
) Node
```
