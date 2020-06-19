package flow_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
)

func TestIdentifier_MarshalJSON(t *testing.T) {
	var tests = []struct {
		name         string
		id           flow.Identifier
		expectedJSON string
	}{
		{
			name:         "Non-empty identifier",
			id:           flow.HexToID("76d3bc41c9f588f7fcd0d5bf4718f8f84b1c41b20882703100b9eb9413807c01"),
			expectedJSON: `"76d3bc41c9f588f7fcd0d5bf4718f8f84b1c41b20882703100b9eb9413807c01"`,
		},
		{
			name:         "Empty identifier",
			id:           flow.EmptyID,
			expectedJSON: `"0000000000000000000000000000000000000000000000000000000000000000"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.id)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedJSON, string(data))

			var id flow.Identifier

			err = json.Unmarshal(data, &id)
			require.NoError(t, err)

			assert.Equal(t, tt.id, id)
		})
	}
}

func TestIdentifier_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name  string
		json  string
		check func(t *testing.T, id flow.Identifier, err error)
	}{
		{
			name: "Non-empty identifier",
			json: `"76d3bc41c9f588f7fcd0d5bf4718f8f84b1c41b20882703100b9eb9413807c01"`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.NoError(t, err)
				assert.Equal(t, flow.HexToID("76d3bc41c9f588f7fcd0d5bf4718f8f84b1c41b20882703100b9eb9413807c01"), id)
			},
		},
		{
			name: "Empty identifier",
			json: `""`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.NoError(t, err)
				assert.Equal(t, flow.EmptyID, id)
			},
		},
		{
			name: "Zero identifier",
			json: `"0000000000000000000000000000000000000000000000000000000000000000"`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.NoError(t, err)
				assert.Equal(t, flow.EmptyID, id)
			},
		},
		{
			name: "Null identifier",
			json: `null`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.NoError(t, err)
				assert.Equal(t, flow.EmptyID, id)
			},
		},
		{
			name: "Too short",
			json: `"76d3bc41c9f588f7fcd0d5bf4718f8f84b1c41b20882703100b9eb941380"`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.Error(t, err)
				assert.Equal(t, flow.EmptyID, id)
			},
		},
		{
			name: "Too long",
			json: `"76d3bc41c9f588f7fcd0d5bf4718f8f84b1c41b20882703100b9eb9413807c01a9b5"`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.Error(t, err)
				assert.Equal(t, flow.EmptyID, id)
			},
		},
		{
			name: "Invalid hex",
			json: `"foobar"`,
			check: func(t *testing.T, id flow.Identifier, err error) {
				assert.Error(t, err)
				assert.Equal(t, flow.EmptyID, id)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id flow.Identifier

			data := []byte(tt.json)
			err := json.Unmarshal(data, &id)

			tt.check(t, id, err)
		})
	}
}
