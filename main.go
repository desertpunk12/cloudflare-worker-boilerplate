package main

// This file exists only to satisfy Go's requirement for a main package.
// When compiling to WASM, the actual entry points are the exported functions
// in wasm.go (renderIndex, renderDynamicContent) that are called from JavaScript.
