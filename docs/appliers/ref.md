---
weight: 1
title: Ref{}
---

```go {style=tokyonight-moon}
type Ref[T any] struct {
	Ptr *T
	Fn  func(T)
}
```
