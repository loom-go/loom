---
title: Signal()
weight: 1
---

```go {style=tokyonight-moon}
func Signal[T any](initial T) (get func() T, set func(T))
```
