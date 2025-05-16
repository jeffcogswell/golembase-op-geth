package entity_test

import (
	"bytes"
	"testing"

	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/entity"
	"github.com/jeffcogswell/golembase-op-geth/rlp"
	"github.com/stretchr/testify/require"
)

func TestActivePayloadRLP(t *testing.T) {
	tests := []struct {
		name    string
		payload entity.EntityMetaData
	}{
		{
			name: "empty payload",
			payload: entity.EntityMetaData{
				ExpiresAtBlock:     0,
				StringAnnotations:  []entity.StringAnnotation{},
				NumericAnnotations: []entity.NumericAnnotation{},
			},
		},
		{
			name: "payload with data",
			payload: entity.EntityMetaData{
				ExpiresAtBlock: 12345,
				StringAnnotations: []entity.StringAnnotation{
					{Key: "key1", Value: "value1"},
					{Key: "key2", Value: "value2"},
				},
				NumericAnnotations: []entity.NumericAnnotation{
					{Key: "num1", Value: 42},
					{Key: "num2", Value: 123},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to RLP
			buf := new(bytes.Buffer)
			err := rlp.Encode(buf, &tt.payload)
			require.NoError(t, err)

			// Unmarshal back from RLP
			var decoded entity.EntityMetaData
			err = rlp.DecodeBytes(buf.Bytes(), &decoded)
			require.NoError(t, err)

			// Verify all fields match
			require.Equal(t, tt.payload.ExpiresAtBlock, decoded.ExpiresAtBlock)
			require.Equal(t, tt.payload.StringAnnotations, decoded.StringAnnotations)
			require.Equal(t, tt.payload.NumericAnnotations, decoded.NumericAnnotations)
		})
	}
}
