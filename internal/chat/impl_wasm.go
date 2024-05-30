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

package chat

import (
	"syscall/js"
)

type chatWasmImpl struct {
	eventCallback func(e Event)
}

func (c *chatWasmImpl) Connect(channel string) {
	js.Global().Set("chatMessage", js.FuncOf(c.chatMessage))
	js.Global().Set("onConnect", js.FuncOf(c.onConnect))

	js.Global().Get("startChat").Invoke(channel)
}

func (c *chatWasmImpl) chatMessage(this js.Value, p []js.Value) interface{} {
	user := p[0].String()
	message := p[1].String()

	c.eventCallback(Event{Type_: Message, Message: message, Sender: user})
	return nil
}

func (c *chatWasmImpl) onConnect(this js.Value, p []js.Value) interface{} {
	c.eventCallback(Event{Type_: Connect})
	return nil
}

func (c *chatWasmImpl) Disconnect() {
}

func (c *chatWasmImpl) OnEvent(callback func(e Event)) {
	c.eventCallback = callback
}

func init() {
	registerImpl(NewWasmImpl)
}

func NewWasmImpl() Chat {
	return &chatWasmImpl{
		eventCallback: func(e Event) {},
	}
}
