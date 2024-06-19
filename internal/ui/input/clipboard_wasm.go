//go:build wasm

/*
 *  Copyright (c) 2024 Juan Medina
 *
 *   All rights reserved. This software and related documentation are proprietary to Juan Medina.
 *
 *   This source code is for internal use only and may not be copied, modified, or distributed
 *   without the express written permission of Juan Medina. Any use of this software for any
 *   purpose other than its intended use is strictly prohibited and may result in severe civil
 *   and criminal penalties.
 *
 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 *   INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
 *   PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
 *   FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 *   OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 *   DEALINGS IN THE SOFTWARE.
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
