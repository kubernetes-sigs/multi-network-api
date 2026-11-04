# Copyright 2026
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: verify

.PHONY: verify
verify:
	hack/verify-all.sh

.PHONY: clean-generated
clean-generated:
	rm -rf pkg/client/clientset
	rm -rf pkg/client/listers
	rm -rf pkg/client/informers

.PHONY: generate
generate: clean-generated update-codegen

.PHONY: update-codegen
update-codegen:
	hack/update-codegen.sh
