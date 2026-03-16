---
title: "YOUR FIRST APP"
weight: 3
---

Firstly, we must create the app itself:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    "fmt"
    "os"

    "github.com/loom-go/loom"
    "github.com/loom-go/term"
)

// define your root component.
// this is what we're going to provide to the renderer.
func App() loom.Node {
    return nil
}

func main() {
    // create a new terminal app
	app := term.NewApp()

    // tell the renderer to render in fullscreen,
    // and give it our root `App` component.
	for err := range app.Run(term.RenderFullscreen, App) {
        // in case of error, close the app and exit with status 1
		app.Close()
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    "github.com/loom-go/loom"
    "github.com/loom-go/web"
)

// define your root component.
// this is what we're going to provide to the renderer.
func App() loom.Node {
    return nil
}

func main() {
    // create a new web app
	app := web.NewApp()

    // tell the renderer to render in `#root`,
    // and give it our root `App` component.
	for err := range app.Run("#root", App) {
        // in case of error, log it to the console
        web.ConsoleError(fmt.Sprintf("Error: %v\n", err)
	}
}
```

{{< /tab >}}
{{< /tabs >}}

With the app setup, let's create our first component:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    . "github.com/loom-go/term/components"
)

func Title() loom.Node {
    return P(Text("beep boop"))
}

// update the root component to use `Title`
func App() loom.Node {
    return Title()
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    . "github.com/loom-go/web/components"
)

func Title() loom.Node {
    return P(Text("beep boop"))
}

// update the root component to use `Title`
func App() loom.Node {
    return Title()
}
```

{{< /tab >}}
{{< /tabs >}}

Now try running the app:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```bash {style=tokyonight-moon}
go run main.go
```

{{< /tab >}}
{{< tab >}}

Create an `index.html` file:

```html {style=tokyonight-moon}
<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>App</title>
  </head>

  <body>
    <div id="app"></div>

    <script src="https://cdn.jsdelivr.net/gh/golang/go@go1.25.0/lib/wasm/wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(
        fetch("main.wasm"),
        go.importObject,
      ).then((r) => go.run(r.instance));
    </script>
  </body>
</html>
```

Build your `main.go` in wasm:

```go {style=tokyonight-moon}
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

And serve your files with your favorite http server!
Here with [`serve`](https://www.npmjs.com/package/serve):

```go {style=tokyonight-moon}
npx serve
```

{{< /tab >}}
{{< /tabs >}}

**✦‧₊˚⋅** And its live! **✦‧₊˚⋅**

### Styling

Styling works via [appliers](/docs/get-started/concepts#applier). How it works will vary depending on the renderer you are using, but it stays mostly the same. Make sure to read your renderer's styling guide to understand more!

Let's try to style our component:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    . "github.com/loom-go/term/components"
)

// create a textStyle variable containing a Style{} applier
var titleStyle = Style{
    Color:      "darkred", // set the text color to darkred
    FontWeight: "Bold", // and the font bold
}

func Title() loom.Node {
    // apply textStyle to the P Node (which will be inherited by the Text Node)
    return P(
        Text("beep boop"),
        Apply(textStyle),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    . "github.com/loom-go/web/components"
)

// create a textStyle variable containing a Style{} applier
var titleStyle = Style{
    "color":       "darkred", // set the text color to darkred
    "font-weight": "bold", // and the font bold
}

func Title() loom.Node {
    // apply textStyle to the P Node (which will be inherited by the Text Node)
    return P(
        Text("beep boop"),
        Apply(textStyle),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

We can go a bit further and make it centered on the screen by styling our `App` component:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
var rootStyle = Style{
    Width:          "100%",
    Height:         "100%",
	PaddingVertical: 6,
    FlexDirection:  "column",
    AlignItems:     "center",
}

func App() loom.Node {
    // wrap our component in a box, and apply our root style on that box
    return Box(
        Title(),
        Apply(rootStyle),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
var rootStyle = Style{
    "width":           "100vw",
    "height":          "100vh",
    "display":         "flex",
    "padding":         "1rem 0",
    "flex-direction":  "column",
    "align-items":     "center",
}

func App() loom.Node {
    // wrap our component in a div, and apply our root style on that box
    return Div(
        Title(),
        Apply(rootStyle),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

If you run the app again, you should see the title is now centered!

### Attributes

Using attributes is very similar to styling. It works with [appliers](/docs/get-started/concepts#appliers) too, and will vary depending on the renderer you're using.

The official renderers each provide an `Attribute{}` applier (and `Attr{}` as an alias).

Let's try adding a new component that uses attributes:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
var titleStyle = Style{
    // ...
    MarginBottom: 3, // add a bit of space for the title
}

var formStyle = Style{
	MaxWidth:      40,
	FlexDirection: "column",
	GapRow:        1,
}

var styleInput = Style{
    Width:                30,
	AlignSelf:            "center",
    BackgroundColor:      "darkgray",
    PlaceholderFontStyle: "italic",
}

func Form() loom.Node {
    return Input(Apply(
        styleInput, // apply our style for this input
        Attr{Placeholder: "Your value..."}, // define a placeholder via attributes
    ))
}

func App() loom.Node {
    return Box(
        Title(),
        Form(), // add our new component to the `App`
        Apply(rootStyle),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
var titleStyle = Style{
    // ...
    "margin-bottom": "1.5rem", // add a bit of space for the title
}

var formStyle = Style{
	"max-width":      "200px",
	"flex-direction": "column",
	"row-gap":        "10px",
}

var styleInput = Style{
    "width":      "100px",
    "align-self": "center",
}

func Form() loom.Node {
    return Box(
        Input(Apply(
            styleInput, // apply our style for this input
            Attr{"placeholder": "Your value..."}, // define a placeholder via attributes
        )),

        Apply(formStyle),
    )
}

func App() loom.Node {
    return Div(
        Title(),
        Form(), // add our new component to the `App`
        Apply(rootStyle),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

Run the app again to see the changes!

### Events

Events is no exception, it also works through appliers.

The official renderers provide the `On{}` applier. It registers callbacks on the Node for spesific events, and removed them when the Node is disposed.

Let's keep going with our `Form()` component by using the `On{}` applier:

{{< tabs items="TERM,WEB" >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    // ...

    "encoding/base64"

    . "github.com/loom-go/loom/components"
)

func Form() loom.Node {
    result, setResult := Signal("")

    onInput := func(evt *term.EventInput) {
        // encode the input's value in base64
        raw := []byte(evt.Value)
        encoded := base64.StdEncoding.EncodeToString(raw)

        // set result with the encoded value
        setResult(encoded)
    }

    return Box(
        Input(Apply(
            styleInput,
            On{Input: onInput}, // register our callback with On{}
            Attr{Placeholder: "Your value..."},
        )),

        // display the encoded result
        P(Text("Result: "), BindText(result)),

        Apply(formStyle),
    )
}
```

{{< /tab >}}
{{< tab >}}

```go {style=tokyonight-moon}
import (
    // ...

    "encoding/base64"

    . "github.com/loom-go/loom/components"
)

func Form() loom.Node {
    result, setResult := Signal("")

    onInput := func(evt *web.EventInput) {
        // encode the input's value in base64
        raw := []byte(evt.Value)
        encoded := base64.StdEncoding.EncodeToString(raw)

        // set result with the encoded value
        setResult(encoded)
    }

    return Div(
        Input(Apply(
            styleInput,
            On{"input": onInput}, // register our callback with On{}
            Attr{"placeholer": "Your value..."},
        )),

        // display the encoded result
        P(Text("Result: "), BindText(result)),

        Apply(formStyle),
    )
}
```

{{< /tab >}}
{{< /tabs >}}

And run the app to see our little base64 encoder!

---

Feel free to keep exploring the docs and the [examples](/docs/examples) to see the various possibilities and features of loom!
