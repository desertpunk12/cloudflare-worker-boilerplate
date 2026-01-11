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
	js.Global().Set("renderHome", js.FuncOf(renderHome))
	js.Global().Set("renderKV", js.FuncOf(renderKV))
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

func renderHome(this js.Value, args []js.Value) interface{} {
	return renderToString(Home())
}

func renderKV(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			defer func() {
				if r := recover(); r != nil {
					reject.Invoke(fmt.Sprintf("Panic in renderKV: %v", r))
				}
			}()

			key := "kv_demo_key"

			// GET
			val, err := KVGet(key)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to get KV value: %v", err))
				return
			}

			displayVal := val
			if displayVal == "" {
				displayVal = "(empty - first run?)"
			}

			// SET
			newVal := fmt.Sprintf("Updated at %s from Go WASM", time.Now().Format(time.RFC1123))
			err = KVSet(key, newVal)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to set KV value: %v", err))
				return
			}

			// Render simple HTML
			html := fmt.Sprintf(`
				<div style="font-family: sans-serif; padding: 2rem; max-width: 600px; margin: 0 auto;">
					<h1 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 0.5rem;">Cloudflare KV + Go WASM</h1>
					
					<div style="background: #f8f9fa; border: 1px solid #e9ecef; border-radius: 8px; padding: 1.5rem; margin-top: 1.5rem;">
						<h3 style="margin-top: 0;">Previous Value:</h3>
						<pre style="background: #e9ecef; padding: 0.5rem; border-radius: 4px;">%s</pre>
					</div>

					<div style="background: #d4edda; color: #155724; border: 1px solid #c3e6cb; border-radius: 8px; padding: 1.5rem; margin-top: 1rem;">
						<strong>Success!</strong> Value has been updated.
						<div style="margin-top: 0.5rem;">New Value: %s</div>
					</div>

					<div style="margin-top: 2rem; text-align: center;">
						<p><small>Refresh the page to see the new value cycle through.</small></p>
						<a href="/" style="color: #3498db; text-decoration: none; font-weight: bold;">&larr; Back to Home</a>
					</div>
				</div>
			`, displayVal, newVal)

			resolve.Invoke(html)
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func renderToString(c templ.Component) string {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return fmt.Sprintf("<div>Error rendering component: %v</div>", err)
	}
	return buf.String()
}
