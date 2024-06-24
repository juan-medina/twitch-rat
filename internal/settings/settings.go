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

const (
	CHANNEL_NAME             = "channel"
	DEFAULT_CHANNEL_NAME     = ""
	SONG_VOLUME              = "song_volume"
	DEFAULT_SONG_VOLUME      = 50
	SOUND_VOLUME             = "sound_volume"
	DEFAULT_SOUND_VOLUME     = 50
	DEBUG                    = "debug"
	DEFAULT_DEBUG            = false
	GAME_MODE                = "game_mode"
	GAME_MODE_AFK            = 0
	GAME_MODE_BATTLE         = 1
	GAME_MODE_CUSTOM         = 2
	JOIN_MODE                = "join_mode"
	JOIN_MODE_CHATTER        = 0
	JOIN_MODE_WITH_COMMAND   = 1
	REJOIN_MODE              = "rejoin_mode"
	REJOIN_MODE_YES          = 0
	REJOIN_MODE_NO           = 1
	ATTACK_MODE              = "attack_mode"
	ATTACK_MODE_RANDOM       = 0
	ATTACK_MODE_WITH_COMMAND = 1
	HEAL_MODE                = "heal_mode"
	HEAL_MODE_RANDOM         = 0
	HEAL_MODE_WITH_COMMAND   = 1
	HEAL_TO                  = "heal_to"
	HEAL_TO_ANYONE           = 0
	HEAL_TO_SELF             = 1
	HEAL_TO_OTHERS           = 2
	RAT_HEALTH               = "rat_health"
	DEFAULT_HEALTH_AFK       = 50
	DEFAULT_HEALTH_BATTLE    = 100
)

type Settings interface {
	Init()

	SetValue(key string, value string)
	GetValue(key string, defaultValue string) string

	GetFloatValue(key string, defaultValue float64) float64
	SetFloatValue(key string, value float64)

	GetBoolValue(key string, defaultValue bool) bool
	SetBoolValue(key string, value bool)

	GetIntValue(key string, defaultValue int) int
	SetIntValue(key string, value int)

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
	if value, err := strconv.ParseFloat(valueStr, 64); err != nil {
		s.settings[key] = strDefaultValue
		return defaultValue
	} else {
		return value
	}
}

func (s *settingsImpl) SetValue(key string, value string) {
	s.settings[key] = value
}

func (s *settingsImpl) SetFloatValue(key string, value float64) {
	s.SetValue(key, strconv.FormatFloat(value, 'f', -1, 64))
}

func (s *settingsImpl) GetBoolValue(key string, defaultValue bool) bool {
	strDefaultValue := strconv.FormatBool(defaultValue)
	valueStr := s.GetValue(key, strDefaultValue)
	if value, err := strconv.ParseBool(valueStr); err != nil {
		s.settings[key] = strDefaultValue
		return defaultValue
	} else {
		return value
	}
}

func (s *settingsImpl) SetBoolValue(key string, value bool) {
	s.SetValue(key, strconv.FormatBool(value))
}

func registerImpl(impl func(application string) Storage) {
	registeredStorage = impl
}

func (s *settingsImpl) SetIntValue(key string, value int) {
	s.SetValue(key, strconv.FormatInt(int64(value), 10))
}

func (s *settingsImpl) GetIntValue(key string, defaultValue int) int {
	strDefaultValue := strconv.FormatInt(int64(defaultValue), 10)
	valueStr := s.GetValue(key, strDefaultValue)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err != nil {
		s.settings[key] = strDefaultValue
		return defaultValue
	} else {
		return int(value)
	}
}

func New(application string) Settings {
	return &settingsImpl{
		application: application,
		settings:    make(map[string]string),
		storage:     registeredStorage(application),
	}
}
