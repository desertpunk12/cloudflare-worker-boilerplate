# Cloudflare Worker Boilerplate (Go/WASM + TinyGo + Templ)

This project is a boilerplate for building **Cloudflare Workers** using **Go** (compiled to WebAssembly via **TinyGo**) and **Templ** for type-safe HTML templating.

It demonstrates how to run Go code on the edge, rendering dynamic HTML content using `a-h/templ`, and serving it via a lightweight Javascript worker shim.

## Features

-   **⚡ TinyGo**: Optimized for small WASM binary sizes and fast startup times.
-   **🧩 Templ**: Type-safe, component-based HTML templating for Go.
-   **🎨 Tailored UI**: Includes setup for modern, responsive designs (with dark mode support).
-   **🚀 Cloudflare Workers**: Deploys globally to the edge.

## Prerequisites

Ensure you have the following installed on your system:

1.  **Go** (1.21+): [Download Go](https://go.dev/dl/)
2.  **TinyGo** (0.30+): [Install TinyGo](https://tinygo.org/getting-started/install/)
    *   *Note: TinyGo is required for the build process defined in the Makefile.*
3.  **Node.js & npm**: [Download Node.js](https://nodejs.org/)
4.  **Wrangler**: The Cloudflare Developer Platform CLI.
    ```bash
    npm install -g wrangler
    ```
5.  **Templ CLI**: To generate Go code from `.templ` files.
    ```bash
    go install github.com/a-h/templ/cmd/templ@latest
    ```

## Setup

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/cloudflare-worker-boilerplate.git
    cd cloudflare-worker-boilerplate
    ```

2.  **Install Go dependencies**:
    ```bash
    go mod download
    ```

## Development

### 1. Generate Templates
If you modify any `.templ` files (e.g., `index.templ`, `dynamic-content.templ`), you must regenerate the Go code:
```bash
templ generate
```

### 2. Build the WASM Module
Compile the Go code into WebAssembly. This step also ensures the correct `wasm_exec.js` glue code is present.
```bash
make build
```

### 3. Run Locally
Start the local Wrangler development server to test your worker:
```bash
wrangler dev
```
*   This will start a local server (usually at `http://localhost:8787`).
*   Press `b` in the terminal to open the browser.

## Deployment

To deploy your worker to the Cloudflare global network:

```bash
make deploy
```
*   This command runs `make build` first, then executes `wrangler deploy`.

## Project Structure

*   **`main.go`, `wasm.go`**: The entry point for the Go WASM application.
*   **`*.templ`**: HTML templates defined using the Templ syntax.
*   **`Makefile`**: Automation instructions for building and deploying.
*   **`worker.js`**: The JavaScript entry point for the Cloudflare Worker. It instantiates the WASM module and passes requests to it.
*   **`wrangler.toml`**: Cloudflare Worker configuration file.

## Troubleshooting

-   **"syscall/js: not supported by TinyGo"**: Ensure you are using `tinygo build` and not standard `go build`.
-   **Wrangler Errors**: Make sure you have authenticated with Cloudflare using `wrangler login`.
