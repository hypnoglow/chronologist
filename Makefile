# Copyright 2018 The Chronologist Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: dep test build

.PHONY: dep
dep:
	dep ensure -v -vendor-only

.PHONY: build
build:
	go build -o ${CURDIR}/bin/chronologist ./cmd/chronologist

.PHONY: test
test:
	go test -v $(shell go list ./... | grep -v "e2e")

.PHONY: test-e2e
test-e2e:
	go test -v ./tests/e2e/...

.PHONY: image
image:
	docker image build -t hypnoglow/chronologist:latest .

. PHONY: mockgen
mockgen:
	minimock -i github.com/hypnoglow/chronologist/internal/grafana.Annotator -o ./internal/grafana/mocks -s _mock.go
