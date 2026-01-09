# Makefile for Cloudflare Worker with TinyGo

# Set environment variables for TinyGo WASM build
export GOOS=js
export GOARCH=wasm

generate:
	@echo "Generating templ templates"
	templ generate


# Default target
.PHONY: build
build: generate
	@echo "Building WASM module with TinyGo..."
	tinygo build -o main.wasm -target wasm .
	@if [ ! -f wasm_exec.js ]; then \
		echo "Copying wasm_exec.js..."; \
		cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js .; \
	else \
		echo "wasm_exec.js already exists. Skipping copy."; \
	fi
	@echo "Build complete."

deploy: build
	@echo Deploying to Cloudflare
	wrangler deploy


.PHONY: clean
clean:
	rm -f main.wasm wasm_exec.js
