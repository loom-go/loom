---
weight: 4
title: For()
---

```go {style=tokyonight-moon}
func For[T comparable](
    items func() []T,
    mapper func(item T, index Accesor[int]) Node,
) Node
```
