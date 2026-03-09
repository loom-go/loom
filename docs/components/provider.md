---
weight: 5
title: Provider()
---

```go {style=tokyonight-moon}
func Provider[T any](ctx Context[T], value T, fn func() Node) Node
```
