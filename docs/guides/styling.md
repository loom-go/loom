---
title: "STYLING"
weight: 3
---

How to style nodes in loom.

This document expects you to be familiar with loom's core concepts and signal-based reactivity. If not -> [CORE CONCEPTS](/docs/concepts), -> [REACTIVITY](/docs/guides/reactivity) and -> [BINDING](/docs/guides/binding)

---

<!--
style priorization:
- on first render, each style property (width, color, etc) is applied top-down. last in the tree wins.
- after that, properties are prioritised winthin each layer.

"blue" wins:
Box(
	Apply(Style{Color: "red"}),
	Apply(Style{Color: "blue"}),
)

"blue" wins:
Box(
	Apply(
		Style{Color: "red"},
		Style{Color: "blue"},
 ),
)

"blue" wins (even when signal changes):
Box(
	Apply(
		Style{Color: getColor},
		Style{Color: "blue"},
 ),
)

"blue" wins, then getColor signal wins on change:
Box(
	Apply(Style{Color: getColor}),
	Apply(Style{Color: "blue"}),
)

"blue" wins, then "red" wins on hover:
Box(
	ApplyOn("hover", Style{Color: "red"}),
	Apply(Style{Color: "blue"}),
)
-->
