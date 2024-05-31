//go:build !wasm

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
	"github.com/gempir/go-twitch-irc/v4"
)

type chatGoTwitchIrcImpl struct {
	client        *twitch.Client
	eventCallback func(e Event)
	alreadyJoined bool
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

	go func() {
		err := c.client.Connect()
		if err != nil {
			panic(err)
		}
	}()
}

func (c *chatGoTwitchIrcImpl) onMessage(message twitch.PrivateMessage) {
	c.eventCallback(Event{Type_: Message, Message: message.Message, Sender: message.User.DisplayName})
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
