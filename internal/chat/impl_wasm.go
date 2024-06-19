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

package chat

import (
	"syscall/js"

	"github.com/juan-medina/twitch-rat/internal/colors"
)

type chatWasmImpl struct {
	eventCallback func(e Event)
	alreadyJoined bool
}

func (c *chatWasmImpl) Connect(channel string) {
	js.Global().Set("chatMessage", js.FuncOf(c.chatMessage))
	js.Global().Set("onSelfJoinMessage", js.FuncOf(c.onSelfJoinMessage))

	js.Global().Get("startChat").Invoke(channel)
}

func (c *chatWasmImpl) chatMessage(this js.Value, p []js.Value) interface{} {
	user := p[0].String()
	message := p[1].String()
	colorStr := p[2].String()

	var userColor = colors.Black
	if colorStr != "" {
		userColor = colors.FromHtml(colorStr)
	}

	c.eventCallback(Event{Type_: Message, Message: message, Sender: user, UserColor: userColor})
	return nil
}

func (c *chatWasmImpl) onSelfJoinMessage(this js.Value, p []js.Value) interface{} {
	if c.alreadyJoined {
		return nil
	}
	c.eventCallback(Event{Type_: Connect})
	c.alreadyJoined = true
	return nil
}

func (c *chatWasmImpl) Disconnect() {
	js.Global().Get("stopChat").Invoke()
	c.alreadyJoined = false
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
