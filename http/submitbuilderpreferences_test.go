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

package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	clienthttp "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/require"
)

func TestSubmitBuilderPreferencesPosts(t *testing.T) {
	tests := []struct {
		name              string
		enforceJSON       bool
		customSpecSupport bool
	}{
		{name: "SSZ"},
		{name: "SSZCustomSpec", customSpecSupport: true},
		{name: "JSON", enforceJSON: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			expected := makeBuilderPreferences(2)
			received := false
			server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				switch r.URL.Path {
				case "/eth/v1/node/version":
					_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
				case "/eth/v1/node/syncing":
					_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
				case "/eth/v1/config/spec":
					_, _ = w.Write([]byte(`{"data":{"MIN_SEED_LOOKAHEAD":"1","SLOTS_PER_EPOCH":"32"}}`))
				case "/eth/v1/validator/builder_preferences":
					received = true
					require.Equal(t, nethttp.MethodPost, r.Method)
					require.Equal(t, "gloas", r.Header.Get("Eth-Consensus-Version"))
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					var decoded []*gloas.BuilderPreferencesEntry
					if test.enforceJSON {
						require.Equal(t, "application/json", r.Header.Get("Content-Type"))
						require.NoError(t, json.Unmarshal(body, &decoded))
					} else {
						require.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
						require.NoError(t, dynssz.GetGlobalDynSsz().UnmarshalSSZ(&decoded, body))
					}
					require.Equal(t, expected, decoded)
					w.WriteHeader(nethttp.StatusOK)
				default:
					w.WriteHeader(nethttp.StatusNotFound)
				}
			}))
			defer server.Close()

			params := []clienthttp.Parameter{
				clienthttp.WithAddress(server.URL),
				clienthttp.WithEnforceJSON(test.enforceJSON),
				clienthttp.WithCustomSpecSupport(test.customSpecSupport),
			}
			service, err := clienthttp.New(ctx, params...)
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).
				SubmitBuilderPreferences(ctx, expected)
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

// TestSubmitBuilderPreferencesReportsRejectedEntries covers a builder rejecting
// its entry: the node's per-entry failures reach the caller in the api.Error.
func TestSubmitBuilderPreferencesReportsRejectedEntries(t *testing.T) {
	ctx := context.Background()
	failures := `{"code":400,"message":"some entries failed","failures":[{"index":0,"message":"builder rejected"}]}`
	received := false
	server := builderPreferencesServer(t, nethttp.StatusBadRequest, failures, &received, "")
	defer server.Close()

	service, err := clienthttp.New(ctx, clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.BuilderPreferencesSubmitter).
		SubmitBuilderPreferences(ctx, makeBuilderPreferences(1))
	require.ErrorContains(t, err, "failed to submit builder preferences")
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, nethttp.StatusBadRequest, apiErr.StatusCode)
	require.JSONEq(t, failures, string(apiErr.Data))
	require.True(t, received)
}

func TestSubmitBuilderPreferencesRequiresExactStatusOK(t *testing.T) {
	ctx := context.Background()
	server := builderPreferencesServer(t, nethttp.StatusNoContent, "", nil, "")
	defer server.Close()

	service, err := clienthttp.New(ctx, clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.BuilderPreferencesSubmitter).
		SubmitBuilderPreferences(ctx, makeBuilderPreferences(1))
	require.EqualError(t, err, "failed to submit builder preferences\nunexpected status code 204")
}

// TestSubmitBuilderPreferencesEnforcesLimit covers the limit of
// MAX_BUILDER_ENTRIES entries per slot of the proposer lookahead, both from the
// mainnet preset and, with custom spec support, from the spec.
func TestSubmitBuilderPreferencesEnforcesLimit(t *testing.T) {
	for _, test := range []struct {
		name  string
		spec  string
		count int
		err   string
	}{
		{name: "StaticAtLimit", count: 4096},
		{name: "StaticOneOverLimit", count: 4097, err: "too many builder preferences"},
		{name: "SpecAtLimit", spec: builderPreferencesSpec(1, 2), count: 256},
		{name: "SpecOneOverLimit", spec: builderPreferencesSpec(1, 2), count: 257, err: "too many builder preferences"},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := builderPreferencesServer(t, nethttp.StatusOK, "", &received, test.spec)
			defer server.Close()
			service, err := clienthttp.New(context.Background(),
				clienthttp.WithAddress(server.URL),
				clienthttp.WithCustomSpecSupport(test.spec != ""),
			)
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).
				SubmitBuilderPreferences(context.Background(), makeBuilderPreferences(test.count))
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				require.ErrorIs(t, err, client.ErrInvalidOptions)
				require.False(t, received)
			} else {
				require.NoError(t, err)
				require.True(t, received)
			}
		})
	}
}

func TestSubmitBuilderPreferencesRejectsInvalidEntries(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(entries []*gloas.BuilderPreferencesEntry)
		err    string
	}{
		{
			name:   "NilEntry",
			mutate: func(entries []*gloas.BuilderPreferencesEntry) { entries[1] = nil },
			err:    "nil builder preference supplied",
		},
		{
			name:   "NilAuth",
			mutate: func(entries []*gloas.BuilderPreferencesEntry) { entries[1].Auth = nil },
			err:    "nil builder preference authorization supplied",
		},
		{
			name:   "NilAuthMessage",
			mutate: func(entries []*gloas.BuilderPreferencesEntry) { entries[1].Auth.Message = nil },
			err:    "nil builder preference authorization supplied",
		},
		{
			name:   "EmptyURL",
			mutate: func(entries []*gloas.BuilderPreferencesEntry) { entries[1].URL = nil },
			err:    "invalid builder preference URL length",
		},
		{
			name: "OversizedURL",
			mutate: func(entries []*gloas.BuilderPreferencesEntry) {
				entries[1].URL = []byte("https://" + strings.Repeat("a", 2041))
			},
			err: "invalid builder preference URL length",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := builderPreferencesServer(t, nethttp.StatusOK, "", &received, "")
			defer server.Close()

			entries := makeBuilderPreferences(2)
			test.mutate(entries)
			service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).
				SubmitBuilderPreferences(context.Background(), entries)
			require.ErrorContains(t, err, test.err)
			require.ErrorIs(t, err, client.ErrInvalidOptions)
			require.False(t, received)
		})
	}
}

func TestSubmitBuilderPreferencesAllowsEmptyList(t *testing.T) {
	for _, test := range []struct {
		name        string
		enforceJSON bool
		body        string
	}{
		{name: "SSZ"},
		{name: "JSON", enforceJSON: true, body: "[]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			received := false
			server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				switch r.URL.Path {
				case "/eth/v1/node/version":
					_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
				case "/eth/v1/node/syncing":
					_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
				case "/eth/v1/validator/builder_preferences":
					received = true
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.Equal(t, test.body, string(body))
					w.WriteHeader(nethttp.StatusOK)
				default:
					w.WriteHeader(nethttp.StatusNotFound)
				}
			}))
			defer server.Close()

			service, err := clienthttp.New(ctx,
				clienthttp.WithAddress(server.URL),
				clienthttp.WithEnforceJSON(test.enforceJSON),
			)
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).
				SubmitBuilderPreferences(ctx, nil)
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

// makeBuilderPreferences returns distinct entries whose URLs and authorization
// data differ in length, so a mis-encoded SSZ offset fails to round trip.
func makeBuilderPreferences(count int) []*gloas.BuilderPreferencesEntry {
	entries := make([]*gloas.BuilderPreferencesEntry, count)
	for i := range entries {
		host := fmt.Sprintf("builder%d.example.com", i)
		entries[i] = &gloas.BuilderPreferencesEntry{
			URL: []byte("https://" + host),
			Auth: &gloas.SignedBuilderRequestAuth{
				Message: &gloas.BuilderRequestAuth{Data: []byte(host), Slot: 100},
			},
			MaxExecutionPayment: 1000000000,
		}
		entries[i].ProposerPubkey[0] = byte(i)
		entries[i].Auth.Signature[0] = byte(i)
	}

	return entries
}

func builderPreferencesSpec(minSeedLookahead, slotsPerEpoch uint64) string {
	return fmt.Sprintf(`{"data":{"MIN_SEED_LOOKAHEAD":"%d","SLOTS_PER_EPOCH":"%d"}}`, minSeedLookahead, slotsPerEpoch)
}

// builderPreferencesServer serves the submission endpoint with the given status
// and response.  An empty spec serves a node with no spec endpoint at all.
func builderPreferencesServer(t *testing.T, status int, response string, received *bool, spec string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		case "/eth/v1/config/spec":
			if spec == "" {
				w.WriteHeader(nethttp.StatusNotFound)

				return
			}
			_, _ = w.Write([]byte(spec))
		case "/eth/v1/validator/builder_preferences":
			if received != nil {
				*received = true
			}
			if response != "" {
				w.Header().Set("Content-Type", "application/json")
			}
			w.WriteHeader(status)
			_, _ = w.Write([]byte(response))
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
}
