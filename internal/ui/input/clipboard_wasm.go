//go:build wasm

/*
 * Copyright (c) 2023 Juan Antonio Medina Iglesias
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package input

import (
	"strconv"
	"syscall/js"
)

type clipboardImplWasm struct {
	callback func(id Id, text string)
}

func (c *clipboardImplWasm) OnReady(callback func(id Id, text string)) {
	c.callback = callback
}

func (c *clipboardImplWasm) Request(id Id) {
	js.Global().Set("onClipboardReady", js.FuncOf(c.onClipboardReady))
	js.Global().Get("requestClipboard").Invoke(strconv.Itoa(int(id)))
}

func (c clipboardImplWasm) onClipboardReady(this js.Value, p []js.Value) interface{} {
	if len(p) == 2 && c.callback != nil {
		idStr := p[0].String()
		text := p[1].String()
		if id, err := strconv.Atoi(idStr); err == nil {
			c.callback(Id(id), text)
		}
	}
	return nil
}

func newClipboardNoWasm() Clipboard {
	return &clipboardImplWasm{}
}

func init() {
	registerClipboardImpl(newClipboardNoWasm)
}
