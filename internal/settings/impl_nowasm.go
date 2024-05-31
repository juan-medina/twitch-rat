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

package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

type noWasmSettings struct {
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

func (n *noWasmSettings) Init(application string) {
	n.application = application
	configPath := configdir.LocalConfig(n.application)
	if err := configdir.MakePath(configPath); err == nil {
		n.settingsFile = filepath.Join(configPath, "settings.json")
		n.load()
	} else {
		panic(err)
	}
}

func (n *noWasmSettings) load() {
	var settings = SettingFile{
		Settings: make([]SettingValue, 0),
	}

	var err error
	if _, err = os.Stat(n.settingsFile); os.IsNotExist(err) {
		fh, err := os.Create(n.settingsFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		encoder := json.NewEncoder(fh)
		encoder.Encode(&settings)
	} else {
		fh, err := os.Open(n.settingsFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		decoder := json.NewDecoder(fh)
		decoder.Decode(&settings)

		for _, setting := range settings.Settings {
			n.settings[setting.Key] = setting.Value
		}
	}
}

func (n *noWasmSettings) Save() {
	var settings = SettingFile{
		Settings: make([]SettingValue, 0),
	}

	for key, value := range n.settings {
		settings.Settings = append(settings.Settings, SettingValue{
			Key:   key,
			Value: value,
		})
	}
	fh, err := os.Create(n.settingsFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	encoder.Encode(&settings)
}

func (n *noWasmSettings) GetValue(key string, defaultValue string) string {
	if value, ok := n.settings[key]; ok {
		return value
	}
	n.settings[key] = defaultValue
	return defaultValue
}

func (n *noWasmSettings) SetValue(key string, value string) {
	n.settings[key] = value
}

func init() {
	registerImpl(NewNoWasmSettings)
}

func NewNoWasmSettings() Settings {
	return &noWasmSettings{
		settings: make(map[string]string),
	}
}
