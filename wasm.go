//go:build js && wasm

package main

import (
	"bytes"
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
)

func main() {
	fmt.Println("Go: main started")
	c := make(chan struct{})
	js.Global().Set("renderIndex", js.FuncOf(renderIndex))
	js.Global().Set("renderDynamicContent", js.FuncOf(renderDynamicContent))
	fmt.Println("Go: exports set, waiting...")
	<-c
}

func renderIndex(this js.Value, args []js.Value) interface{} {
	// args[0] might be the title or other config if we wanted
	// For now, hardcode or use default data, or parse args JSON
	title := "Cloudflare Worker + Go + Templ"
	if len(args) > 0 && args[0].Type() == js.TypeString {
		title = args[0].String()
	}

	component := Index(
		title,
		"Welcome to Cloudflare Workers",
		"Powered by Go, WebAssembly, and templ",
		"Static Items",
		[]string{"Fast", "Secure", "Scalable", "Go-powered"},
		"Load Dynamic Content",
		"/dynamic",
	)

	return renderToString(component)
}

func renderDynamicContent(this js.Value, args []js.Value) interface{} {
	now := time.Now()
	items := []string{
		fmt.Sprintf("Item generated at %s", now.Format(time.TimeOnly)),
		"Another dynamic item",
		"Random Value: " + fmt.Sprint(now.UnixNano()),
	}

	component := DynamicContent(
		"Dynamic Data",
		items,
		now.Format(time.RFC3339),
		now.Format(time.RFC1123),
		now.Format(time.Kitchen),
		"/dynamic",
	)

	return renderToString(component)
}

func renderToString(c templ.Component) string {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return fmt.Sprintf("<div>Error rendering component: %v</div>", err)
	}
	return buf.String()
}
