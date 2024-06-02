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
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

type noWasmStorage struct {
	application  string
	settingsFile string
}

func (n *noWasmStorage) Load() string {
	var data string
	var err error
	if _, err = os.Stat(n.settingsFile); os.IsNotExist(err) {
		fh, err := os.Create(n.settingsFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

	} else {
		if bytes, err := os.ReadFile(n.settingsFile); err == nil {
			data = string(bytes)
		}
	}
	return data
}

func (n *noWasmStorage) Save(data string) {
	fh, err := os.Create(n.settingsFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	fh.WriteString(data)
}

func init() {
	registerImpl(NewNoWasmStorage)
}

func NewNoWasmStorage(application string) Storage {
	var settingsFile string
	configPath := configdir.LocalConfig(application)
	if err := configdir.MakePath(configPath); err == nil {
		settingsFile = filepath.Join(configPath, "settings.json")
	} else {
		panic(err)
	}

	return &noWasmStorage{
		application:  application,
		settingsFile: settingsFile,
	}
}
