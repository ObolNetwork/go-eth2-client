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

package multi_test

import (
	"context"
	"errors"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// erroringBuilderPreferencesClient rejects every submission with the given error.
type erroringBuilderPreferencesClient struct {
	name  string
	err   error
	calls int
}

func (c *erroringBuilderPreferencesClient) Name() string    { return c.name }
func (c *erroringBuilderPreferencesClient) Address() string { return c.name }
func (*erroringBuilderPreferencesClient) IsActive() bool    { return true }
func (*erroringBuilderPreferencesClient) IsSynced() bool    { return true }

func (c *erroringBuilderPreferencesClient) SubmitBuilderPreferences(_ context.Context,
	_ *api.SubmitBuilderPreferencesOpts,
) error {
	c.calls++

	return c.err
}

func builderPreferencesOpts() *api.SubmitBuilderPreferencesOpts {
	return &api.SubmitBuilderPreferencesOpts{
		Preferences: []*gloas.BuilderPreferencesEntry{{
			URL:  []byte("https://builder.example.com"),
			Auth: &gloas.SignedBuilderRequestAuth{Message: &gloas.BuilderRequestAuth{}},
		}},
	}
}

func TestSubmitBuilderPreferences(t *testing.T) {
	ctx := context.Background()

	client1, err := mock.New(ctx, mock.WithName("mock 1"))
	require.NoError(t, err)
	client2, err := mock.New(ctx, mock.WithName("mock 2"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client1, client2}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.BuilderPreferencesSubmitter).
		SubmitBuilderPreferences(ctx, builderPreferencesOpts())
	require.NoError(t, err)
}

// TestSubmitBuilderPreferencesDoesNotFailOver ensures that neither invalid
// options nor a builder rejecting its entry, which the node reports as status
// 400 after submitting the other entries, fails over or deactivates the client.
func TestSubmitBuilderPreferencesDoesNotFailOver(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
	}{
		{
			name: "InvalidOptions",
			err:  errors.Join(errors.New("too many builder preferences"), consensusclient.ErrInvalidOptions),
		},
		{
			name: "BuilderRejected",
			err:  &api.Error{Method: "POST", Endpoint: "/eth/v1/validator/builder_preferences", StatusCode: 400},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			client1 := &erroringBuilderPreferencesClient{name: "erroring 1", err: test.err}
			client2 := &erroringBuilderPreferencesClient{name: "erroring 2", err: test.err}

			multiClient, err := multi.New(ctx,
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]consensusclient.Service{client1, client2}),
			)
			require.NoError(t, err)

			err = multiClient.(consensusclient.BuilderPreferencesSubmitter).
				SubmitBuilderPreferences(ctx, builderPreferencesOpts())
			require.ErrorIs(t, err, test.err)
			require.Equal(t, 1, client1.calls)
			require.Equal(t, 0, client2.calls)
			require.Equal(t, "erroring 1", multiClient.Address())
		})
	}
}
