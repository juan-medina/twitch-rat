# Copyright (c) 2023 Juan Antonio Medina Iglesias
#
#  Permission is hereby granted, free of charge, to any person obtaining a copy
#  of this software and associated documentation files (the "Software"), to deal
#  in the Software without restriction, including without limitation the rights
#  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#  copies of the Software, and to permit persons to whom the Software is
#  furnished to do so, subject to the following conditions:
#
#  The above copyright notice and this permission notice shall be included in
#  all copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
#  THE SOFTWARE.

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
	copy web/*.* build/web/
	copy web/js/*.* build/web/js/
endif
#set GOOS & GOARCH to wasm
	$(GOCMD) env -w GOOS=js GOARCH=wasm
	$(GOBUILD) -o build/web/twitch_rat.wasm $(APP_PATH)	
#restore the original GOOS & GOARCH
	$(GOCMD) env -w GOOS=$(GOOS) GOARCH=$(GOARCH)

web: wasm
	$(info "running web on http://localhost:8000/")
	python -m http.server --directory build/web
