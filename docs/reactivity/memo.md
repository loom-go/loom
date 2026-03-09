---
title: Memo()
weight: 2
---

```go {style=tokyonight-moon}
func Memo[T any](compute func() T) (get func() T)
```
