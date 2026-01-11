//go:build js && wasm

package main

import (
	"errors"
	"syscall/js"
)

// KVGet gets a value from the KV namespace binding attached to globalThis.KV
func KVGet(key string) (string, error) {
	kv := js.Global().Get("KV")
	if kv.IsUndefined() {
		return "", errors.New("KV binding not found on global scope")
	}

	// kv.get(key) returns a Promise
	promise := kv.Call("get", key)
	result, err := await(promise)
	if err != nil {
		return "", err
	}

	if result.IsNull() || result.IsUndefined() {
		return "", nil // Key not found
	}

	return result.String(), nil
}

// KVSet sets a value in the KV namespace
func KVSet(key, value string) error {
	kv := js.Global().Get("KV")
	if kv.IsUndefined() {
		return errors.New("KV binding not found on global scope")
	}

	// kv.put(key, value) returns a Promise
	promise := kv.Call("put", key, value)
	_, err := await(promise)
	return err
}

// await waits for a JS promise to resolve or reject
// It relies on the Go scheduler yielding to the JS event loop while waiting on the channel.
func await(promise js.Value) (js.Value, error) {
	resultCh := make(chan js.Value)
	errCh := make(chan error)

	then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Promise resolved
		var res js.Value
		if len(args) > 0 {
			res = args[0]
		}
		resultCh <- res
		return nil
	})
	defer then.Release()

	catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Promise rejected
		errStr := "unknown javascript error"
		if len(args) > 0 {
			errStr = args[0].String()
			// If it's an Error object, try to get .message
			if args[0].Type() == js.TypeObject && !args[0].Get("message").IsUndefined() {
				errStr = args[0].Get("message").String()
			}
		}
		errCh <- errors.New(errStr)
		return nil
	})
	defer catch.Release()

	promise.Call("then", then).Call("catch", catch)

	select {
	case res := <-resultCh:
		return res, nil
	case err := <-errCh:
		return js.Undefined(), err
	}
}
