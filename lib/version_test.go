/*
Copyright © 2024 Kaley Main <kaleymain@google.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package lib

import (
	"testing"
)

func TestMaxVersion(t *testing.T) {
	got := maxVersion([]string{"", "  ", "not a version"})
	if got != "" {
		t.Errorf("maxVersion({\"\", \"  \", \"not a version\"}) = %v; want empty string", got)
	}

	got = maxVersion([]string{"v1", "v0.3", "v1.2.5", "   "})
	if got != "v1.2.5" {
		t.Errorf(
			"maxVersion({\"v1\", \"v0.3\", \"v1.2.5\", \"   \"}) = %v; want v1.2.5", got)
	}
}
