# Copyright (c) 2024 Juan Medina
# 
#  All rights reserved. This software and related documentation are proprietary to Juan Medina.
# 
#  This source code is for internal use only and may not be copied, modified, or distributed 
#  without the express written permission of Juan Medina. Any use of this software for any 
#  purpose other than its intended use is strictly prohibited and may result in severe civil 
#  and criminal penalties.
# 
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
#  INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR 
#  PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE 
#  FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR 
#  OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
#  DEALINGS IN THE SOFTWARE.


GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTOOL=$(GOCMD) tool
GOFORMAT=$(GOCMD) fmt
GORUN=$(GOCMD) run
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
BUILD_DIR=build

ifeq ($(OS),Windows_NT)
	BINARY_NAME=$(BUILD_DIR)/desktop/twitch_rat.exe
else
	BINARY_NAME=$(BUILD_DIR)/desktop/twitch_rat
endif	

APP_PATH="./internal/app"

# get current GOROOT into a variable
GOROOT=$(shell $(GOCMD) env GOROOT)

#save current GOOS into a variable
GOOS=$(shell $(GOCMD) env GOOS)
#save current GOARCH into a variable
GOARCH=$(shell $(GOCMD) env GOARCH)

default: build

build: clean
	$(GOBUILD) -o $(BINARY_NAME) -v $(APP_PATH)
vet:
	$(GOVET) "./internal/..."
clean:
	$(GOCLEAN) $(APP_PATH)
format:
	$(GOFORMAT) "./internal/..."
tidy:
	$(GOMOD) tidy
run: build
	./$(BINARY_NAME)
update:
	$(GOGET) -u all
	$(GOMOD) tidy
wasm:
ifeq ($(OS),Windows_NT)
	if not exist build mkdir build
	if not exist build\web mkdir build\web
	if not exist build\web\js mkdir build\web\js
	copy $(GOROOT)\misc\wasm\wasm_exec.js build\web\wasm_exec.js
	copy web\*.* build\web\\
	copy web\js\*.* build\web\js\\
else
	mkdir -p build/web/js
	cp $(GOROOT)/misc/wasm/wasm_exec.js build/web/wasm_exec.js
	cp web/*.* build/web/
	cp web/js/*.* build/web/js/
endif
#set GOOS & GOARCH to wasm
	$(GOCMD) env -w GOOS=js GOARCH=wasm
	$(GOBUILD) -o build/web/twitch_rat.wasm $(APP_PATH)	
#restore the original GOOS & GOARCH
	$(GOCMD) env -w GOOS=$(GOOS) GOARCH=$(GOARCH)

web: wasm
	$(info "running web on http://localhost:8000/")
	python -m http.server --directory build/web
publish:
# give error if parameter VERSION has not been provided
ifndef VERSION
	$(error "VERSION is not set")
else
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push --tags
	gh release create v$(VERSION) --generate-notes
endif
