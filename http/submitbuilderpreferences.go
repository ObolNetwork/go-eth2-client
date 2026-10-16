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
	"bytes"
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
)

// maxBuilderEntries is MAX_BUILDER_ENTRIES, the number of builders a proposer may
// configure.
const maxBuilderEntries uint64 = 64

// maxBuilderURLSize is MAX_BUILDER_URL_SIZE.
const maxBuilderURLSize = 2048

// SubmitBuilderPreferences submits builder preferences, which the beacon node
// forwards to each entry's builder.  If any builder rejects its entry the node
// still submits the others and returns status 400, with the per-entry failures
// in the returned api.Error's data.
func (s *Service) SubmitBuilderPreferences(ctx context.Context, opts *api.SubmitBuilderPreferencesOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	if len(opts.Preferences) > 0 &&
		uint64(len(opts.Preferences)) > maxBuilderEntries*s.proposerPreferencesLimit(ctx) {
		return errors.Join(errors.New("too many builder preferences"), client.ErrInvalidOptions)
	}
	for _, entry := range opts.Preferences {
		if entry == nil {
			return errors.Join(errors.New("nil builder preference supplied"), client.ErrInvalidOptions)
		}
		// SSZ marshaling substitutes a zero value for a nil authorization, so without
		// this check the same input would silently submit a zeroed authorization over
		// SSZ and be rejected by the node as a null one over JSON.
		if entry.Auth == nil || entry.Auth.Message == nil {
			return errors.Join(errors.New("nil builder preference authorization supplied"), client.ErrInvalidOptions)
		}
		if len(entry.URL) == 0 || len(entry.URL) > maxBuilderURLSize {
			return errors.Join(errors.New("invalid builder preference URL length"), client.ErrInvalidOptions)
		}
	}

	requestPreferences := apiv1gloas.BuilderPreferencesEntryList(opts.Preferences)
	if requestPreferences == nil {
		requestPreferences = apiv1gloas.BuilderPreferencesEntryList{}
	}
	// The generated list marshaller has a pointer receiver.
	body, contentType, err := s.marshalRequestBody(ctx, &requestPreferences)
	if err != nil {
		return err
	}

	response, err := s.post(ctx,
		"/eth/v1/validator/builder_preferences",
		"",
		&opts.Common,
		bytes.NewReader(body),
		contentType,
		map[string]string{"Eth-Consensus-Version": strings.ToLower(spec.DataVersionGloas.String())},
	)
	if err != nil {
		return errors.Join(errors.New("failed to submit builder preferences"), err)
	}
	if response.statusCode != nethttp.StatusOK {
		return errors.Join(
			errors.New("failed to submit builder preferences"),
			fmt.Errorf("unexpected status code %d", response.statusCode),
		)
	}

	return nil
}
