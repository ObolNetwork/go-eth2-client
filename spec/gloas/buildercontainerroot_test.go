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

package gloas_test

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// The builder API types are regular SSZ containers, unlike the progressive
// containers of the Gloas consensus types.  These roots are calculated by hand
// from the SSZ merkleization rules, independently of the generated code.

func hashPair(a, b [32]byte) [32]byte {
	return sha256.Sum256(append(a[:], b[:]...))
}

// merkleize pads the chunks with zero chunks to the next power of two of the limit
// and returns the root of the resulting tree.
func merkleize(chunks [][32]byte, limit int) [32]byte {
	width := 1
	for width < limit {
		width *= 2
	}
	layer := make([][32]byte, width)
	copy(layer, chunks)
	for len(layer) > 1 {
		next := make([][32]byte, len(layer)/2)
		for i := range next {
			next[i] = hashPair(layer[2*i], layer[2*i+1])
		}
		layer = next
	}

	return layer[0]
}

func packBytes(data []byte) [][32]byte {
	chunks := make([][32]byte, (len(data)+31)/32)
	for i := range chunks {
		copy(chunks[i][:], data[32*i:])
	}

	return chunks
}

func uint64Root(v uint64) [32]byte {
	var chunk [32]byte
	binary.LittleEndian.PutUint64(chunk[:], v)

	return chunk
}

func byteListRoot(data []byte, maxLen int) [32]byte {
	return hashPair(merkleize(packBytes(data), (maxLen+31)/32), uint64Root(uint64(len(data))))
}

func byteVectorRoot(data []byte) [32]byte {
	chunks := packBytes(data)

	return merkleize(chunks, len(chunks))
}

func TestBuilderRequestAuthIsRegularContainer(t *testing.T) {
	auth := &gloas.SignedBuilderRequestAuth{
		Message: &gloas.BuilderRequestAuth{
			Data: []byte("builder.example.com"),
			Slot: 12345,
		},
		Signature: phase0.BLSSignature{0x01, 0x02, 0x03},
	}

	messageRoot := merkleize([][32]byte{
		byteListRoot(auth.Message.Data, 4096),
		uint64Root(uint64(auth.Message.Slot)),
	}, 2)
	root, err := auth.Message.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, messageRoot, root)

	signedRoot := merkleize([][32]byte{messageRoot, byteVectorRoot(auth.Signature[:])}, 2)
	root, err = auth.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, signedRoot, root)
}

func TestBuilderPreferencesEntryIsRegularContainer(t *testing.T) {
	entry := &gloas.BuilderPreferencesEntry{
		ProposerPubkey: phase0.BLSPubKey{0xaa, 0xbb},
		URL:            []byte("https://builder.example.com"),
		Auth: &gloas.SignedBuilderRequestAuth{
			Message:   &gloas.BuilderRequestAuth{Data: []byte("builder.example.com"), Slot: 7},
			Signature: phase0.BLSSignature{0x04},
		},
		MaxExecutionPayment: 1000000000,
	}

	authRoot, err := entry.Auth.HashTreeRoot()
	require.NoError(t, err)
	expected := merkleize([][32]byte{
		byteVectorRoot(entry.ProposerPubkey[:]),
		byteListRoot(entry.URL, 2048),
		authRoot,
		uint64Root(uint64(entry.MaxExecutionPayment)),
	}, 4)

	root, err := entry.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, expected, root)
}
