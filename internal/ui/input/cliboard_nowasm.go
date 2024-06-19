//go:build !wasm

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
	"github.com/f1bonacc1/glippy"
)

type clipboardImplNoWasm struct {
	callback func(id Id, text string)
}

func (c *clipboardImplNoWasm) OnReady(callback func(id Id, text string)) {
	c.callback = callback
}

func (c *clipboardImplNoWasm) Request(id Id) {
	text, _ := glippy.Get()
	if text != "" && c.callback != nil {
		c.callback(id, text)
	}
}

func newClipboardNoWasm() Clipboard {
	return &clipboardImplNoWasm{}
}

func init() {
	registerClipboardImpl(newClipboardNoWasm)
}
