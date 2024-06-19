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

package settings

import (
	"encoding/json"
	"strconv"
)

type Settings interface {
	Init()

	SetValue(key string, value string)
	GetValue(key string, defaultValue string) string

	GetFloatValue(key string, defaultValue float64) float64
	SetFloatValue(key string, value float64)

	GetBoolValue(key string, defaultValue bool) bool
	SetBoolValue(key string, value bool)

	Save()
}

type Storage interface {
	Load() string
	Save(data string)
}

var registeredStorage func(string) Storage = func(application string) Storage {
	panic("settings implementation not registered")
}

type settingsImpl struct {
	application string
	settings    map[string]string
	storage     Storage
}

type jsonSettingValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type jsonSettings struct {
	Settings []jsonSettingValue `json:"settings"`
}

func (s *settingsImpl) Init() {
	data := s.storage.Load()

	var settings = jsonSettings{
		Settings: make([]jsonSettingValue, 0),
	}

	if err := json.Unmarshal([]byte(data), &settings); err == nil {
		for _, setting := range settings.Settings {
			s.settings[setting.Key] = setting.Value
		}
	}
}

func (s *settingsImpl) Save() {
	var settings = jsonSettings{
		Settings: make([]jsonSettingValue, 0),
	}

	for key, value := range s.settings {
		settings.Settings = append(settings.Settings, jsonSettingValue{
			Key:   key,
			Value: value,
		})
	}

	if data, err := json.Marshal(settings); err == nil {
		s.storage.Save(string(data))
	}
}

func (s *settingsImpl) GetValue(key string, defaultValue string) string {
	if value, ok := s.settings[key]; ok {
		return value
	} else {
		s.settings[key] = defaultValue
		return defaultValue
	}
}

func (s *settingsImpl) GetFloatValue(key string, defaultValue float64) float64 {
	strDefaultValue := strconv.FormatFloat(defaultValue, 'f', -1, 64)
	valueStr := s.GetValue(key, strDefaultValue)
	value, _ := strconv.ParseFloat(valueStr, 64)
	return value

}

func (s *settingsImpl) SetValue(key string, value string) {
	s.settings[key] = value
}

func (s *settingsImpl) SetFloatValue(key string, value float64) {
	s.SetValue(key, strconv.FormatFloat(value, 'f', -1, 64))
}

func (s settingsImpl) GetBoolValue(key string, defaultValue bool) bool {
	strDefaultValue := strconv.FormatBool(defaultValue)
	valueStr := s.GetValue(key, strDefaultValue)
	value, _ := strconv.ParseBool(valueStr)
	return value
}

func (s *settingsImpl) SetBoolValue(key string, value bool) {
	s.SetValue(key, strconv.FormatBool(value))
}

func registerImpl(impl func(application string) Storage) {
	registeredStorage = impl
}

func New(application string) Settings {
	return &settingsImpl{
		application: application,
		settings:    make(map[string]string),
		storage:     registeredStorage(application),
	}
}
