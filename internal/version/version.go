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

package version

import (
	"embed"
	"strconv"
	"strings"

	"github.com/juan-medina/twitch-rat/internal/colors"
)

const (
	VERSION_REMOTE = "https://juan-medina.com/twitch-rat/version.txt"
	VERSION_LOCAL  = "embed/version.txt"
)

type Version interface {
	Init()
	Current() versionInfo
	Latest() versionInfo
	Outdated() bool
}

type versionInfo struct {
	Text   string
	Bbcode string
	Major  int
	Minor  int
	Patch  int
	Build  int
}
type versionImpl struct {
	current    versionInfo
	latest     versionInfo
	fileSystem embed.FS
}

func (v versionImpl) Current() versionInfo {
	return v.current
}

func (v *versionImpl) Latest() versionInfo {
	return v.latest
}

func (v *versionImpl) Init() {
	v.fetchCurrent(VERSION_LOCAL)
	v.fetchLatest(VERSION_REMOTE)
}

func (v *versionImpl) info(versionStr string) versionInfo {
	parts := strings.Split(versionStr, ".")
	bbcode := colors.Blue.BBCoded("v") +
		colors.Green.BBCoded(parts[0]) + "." +
		colors.Yellow.BBCoded(parts[1]) + "." +
		colors.Orange.BBCoded(parts[2]) + "." +
		colors.Red.BBCoded(parts[3])
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	build, _ := strconv.Atoi(parts[3])
	return versionInfo{
		Text:   versionStr,
		Bbcode: bbcode,
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Build:  build,
	}
}

func (v *versionImpl) fetchCurrent(filePath string) {
	versionStr := "0.0.0.0"
	if data, err := v.fileSystem.ReadFile(filePath); err == nil {
		versionStr = string(data)
	} else {
		panic(err)
	}
	v.current = v.info(versionStr)
}

func (v *versionImpl) updateLatestCallback(versionStr string) {
	v.latest = v.info(versionStr)
}

func (v versionImpl) Outdated() bool {
	if v.current.Major < v.latest.Major {
		return true
	}
	if v.current.Minor < v.latest.Minor {
		return true
	}
	if v.current.Patch < v.latest.Patch {
		return true
	}
	if v.current.Build < v.latest.Build {
		return true
	}
	return false
}

func New(fileSystem embed.FS) Version {
	return &versionImpl{
		fileSystem: fileSystem,
	}
}
