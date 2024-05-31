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

package settings

import (
	"encoding/base64"
	"encoding/json"
	"syscall/js"
)

type wasmSettings struct {
	application  string
	settings     map[string]string
	settingsFile string
}

type SettingValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SettingFile struct {
	Settings []SettingValue `json:"settings"`
}

func (n *wasmSettings) Init(application string) {
	n.application = application
	n.load()
}

func (n *wasmSettings) load() {
	result := js.Global().Get("getSettings").Invoke(n.application).String()
	if result == "" {
		return
	}

	data, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return
	}

	var settings = SettingFile{
		Settings: make([]SettingValue, 0),
	}
	if err := json.Unmarshal(data, &settings); err == nil {
		for _, setting := range settings.Settings {
			n.settings[setting.Key] = setting.Value
		}
	}
}

func (n *wasmSettings) Save() {
	var settings = SettingFile{
		Settings: make([]SettingValue, 0),
	}

	for key, value := range n.settings {
		settings.Settings = append(settings.Settings, SettingValue{
			Key:   key,
			Value: value,
		})
	}

	data, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	encodedData := []byte(base64.StdEncoding.EncodeToString([]byte(string(data))))
	js.Global().Get("setSettings").Invoke(n.application, string(encodedData))
}

func (n *wasmSettings) GetValue(key string, defaultValue string) string {
	if value, ok := n.settings[key]; ok {
		return value
	}
	n.settings[key] = defaultValue
	return defaultValue
}

func (n *wasmSettings) SetValue(key string, value string) {
	n.settings[key] = value
}

func init() {
	registerImpl(NewWasmSettings)
}

func NewWasmSettings() Settings {
	return &wasmSettings{
		settings: make(map[string]string),
	}
}
