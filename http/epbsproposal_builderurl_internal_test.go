// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSameBuilderURL pins the RFC 3986 section 6.2 normalisation builder URLs are compared under,
// which is also what beacon nodes apply when they parse and echo a URL.
func TestSameBuilderURL(t *testing.T) {
	tests := []struct {
		a, b string
		same bool
	}{
		// Empty path is "/" (6.2.3).
		{a: "https://builder.example", b: "https://builder.example/", same: true},
		// Scheme and host case (6.2.2.1).
		{a: "HTTPS://Builder.Example/", b: "https://builder.example", same: true},
		// Default port (6.2.3).
		{a: "https://builder.example:443", b: "https://builder.example/", same: true},
		{a: "http://builder.example:80/bids", b: "http://builder.example/bids", same: true},
		// Dot segments (6.2.2.3).
		{a: "https://builder.example/a/./b", b: "https://builder.example/a/b", same: true},
		{a: "https://builder.example/a/../b", b: "https://builder.example/b", same: true},
		// A trailing slash of a non-empty path names another resource.
		{a: "https://builder.example/bids", b: "https://builder.example/bids/", same: false},
		{a: "https://builder.example:8443", b: "https://builder.example", same: false},
		{a: "http://builder.example", b: "https://builder.example", same: false},
		{a: "https://builder.example/?a=1", b: "https://builder.example/?a=2", same: false},
		{a: "https://user:pass@builder.example", b: "https://builder.example", same: false},
		{a: "https://builder.example/a", b: "https://builder.example/A", same: false},
		// Unparsable URLs only match themselves.
		{a: "https://builder.example/%zz", b: "https://builder.example/%zz", same: true},
		{a: "https://builder.example/%zz", b: "https://builder.example/%ZZ", same: false},
	}
	for _, test := range tests {
		t.Run(test.a+" "+test.b, func(t *testing.T) {
			require.Equal(t, test.same, sameBuilderURL(test.a, test.b))
			require.Equal(t, test.same, sameBuilderURL(test.b, test.a))
		})
	}
}
