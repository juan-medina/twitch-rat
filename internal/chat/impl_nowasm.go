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

package chat

import (
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/juan-medina/twitch-rat/internal/colors"
)

type chatGoTwitchIrcImpl struct {
	client        *twitch.Client
	eventCallback func(e Event)
	alreadyJoined bool
}

func (c *chatGoTwitchIrcImpl) connectionGoRoutine() {
	if err := c.client.Connect(); err != nil {
		if err != twitch.ErrClientDisconnected {
			panic(err)
		}
	}
	c.eventCallback(Event{Type_: Disconnect})
}

// Connect implements chatBackend.
func (c *chatGoTwitchIrcImpl) Connect(channel string) {
	if c.client != nil {
		panic("We already have a twitch client")
	}
	c.client = twitch.NewAnonymousClient()
	c.client.Join(channel)
	c.client.OnPrivateMessage(c.onMessage)
	c.client.OnSelfJoinMessage(c.onSelfJoinMessage)

	go c.connectionGoRoutine()
}

func (c *chatGoTwitchIrcImpl) onMessage(message twitch.PrivateMessage) {
	var userColor = colors.Black
	if colorStr, ok := message.Tags["color"]; ok {
		userColor = colors.FromHtml(colorStr)
	}
	c.eventCallback(Event{Type_: Message, Message: message.Message, Sender: message.User.DisplayName, UserColor: userColor})
}

func (c *chatGoTwitchIrcImpl) onSelfJoinMessage(twitch.UserJoinMessage) {
	if c.alreadyJoined {
		return
	}
	c.eventCallback(Event{Type_: Connect})
	c.alreadyJoined = true
}

func (c *chatGoTwitchIrcImpl) Disconnect() {
	c.client.Disconnect()
	c.client = nil
	c.alreadyJoined = false
}

// OnEvent implements Chat.
func (c *chatGoTwitchIrcImpl) OnEvent(callback func(e Event)) {
	c.eventCallback = callback
}

func init() {
	registerImpl(NewGoTwitchImpl)
}

func NewGoTwitchImpl() Chat {
	return &chatGoTwitchIrcImpl{
		client:        nil,
		eventCallback: func(e Event) {},
	}
}
