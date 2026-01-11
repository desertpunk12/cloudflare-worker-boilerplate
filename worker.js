import "./wasm_exec.js";

const go = new Go();
let wasmInstance;

async function initWasm(env) {
  if (!wasmInstance) {
    // In Cloudflare Workers, WASM modules are often bound as variables or imported
    // Here we assume 'main.wasm' is handled by the bundler/wrangler and available.
    // However, Wrangler 2+ usually imports wasm as a module if configured or we use the import.
    // Let's rely on the module import standard for workers if possible, or fetch it.

    // BUT the standard 'wasm_exec.js' expects to instantiate streaming or from array buffer.
    // We will import the wasm module directly which Wrangler supports.

    // Note: The below import assumes we are using Esbuild/Webpack/Wrangler's module support
    // which effectively makes the .wasm file available.
    // Usually for 'CompiledWasm' rule, we can import it.

    // Let's try to import the wasm file (wrangler handles this)
    const wasm = await import("./main.wasm");

    // 'wasm' export is usually the Module object or we need to instantiate it.
    // Wrangler "rules" -> CompiledWasm means `import mod from './main.wasm'` gives a WebAssembly.Module.

    wasmInstance = await WebAssembly.instantiate(wasm.default, go.importObject);
    // Monitor Go execution
    go.run(wasmInstance);

    // Wait for Go to initialize exports
    let retries = 0;
    while (typeof globalThis.renderIndex !== "function" && retries < 40) {
      await new Promise((r) => setTimeout(r, 50));
      retries++;
    }

    if (typeof globalThis.renderIndex !== "function") {
      console.error("WASM failed to initialize exports in time.");
    }
  }
}

export default {
  async fetch(request, env, ctx) {
    try {
      await initWasm(env);

      // Try to serve static assets first
      if (env.ASSETS) {
        try {
          const assetResponse = await env.ASSETS.fetch(request);
          if (assetResponse.status !== 404) {
            return assetResponse;
          }
        } catch (e) {
          // Fall through to dynamic routes
        }
      }

      const url = new URL(request.url);

      if (url.pathname === "/dynamic") {
        if (typeof globalThis.renderDynamicContent !== "function") {
          throw new Error(
            "renderDynamicContent is not defined. WASM may not have initialized correctly.",
          );
        }
        const html = globalThis.renderDynamicContent();
        return new Response(html, {
          headers: { "Content-Type": "text/html" },
        });
      }

      if (url.pathname === "/home") {
        if (typeof globalThis.renderHome !== "function") {
          throw new Error(
            "renderHome is not defined. WASM may not have initialized correctly.",
          );
        }
        const html = globalThis.renderHome();
        return new Response(html, {
          headers: { "Content-Type": "text/html" },
        });
      }

      if (url.pathname === "/kv") {
        // Bind the KV namespace to global scope so Go can find it using js.Global().Get("KV")
        // 'KV_DEMO' must match the binding name in wrangler.toml
        globalThis.KV = env.KV_DEMO;

        if (typeof globalThis.renderKV !== "function") {
          throw new Error(
            "renderKV is not defined. WASM may not have initialized correctly.",
          );
        }

        // renderKV returns a Promise, so we await it
        const html = await globalThis.renderKV();
        return new Response(html, {
          headers: { "Content-Type": "text/html" },
        });
      }

      if (url.pathname === "/") {
        if (typeof globalThis.renderIndex !== "function") {
          throw new Error(
            "renderIndex is not defined. WASM may not have initialized correctly.",
          );
        }
        const html = globalThis.renderIndex("Wingo Cloduflare Worker App");
        return new Response(html, {
          headers: { "Content-Type": "text/html" },
        });
      }

      return new Response("Not Found", { status: 404 });
    } catch (err) {
      console.error(err);
      return new Response(`Error: ${err.message}\nStack: ${err.stack}`, {
        status: 500,
      });
    }
  },
};
