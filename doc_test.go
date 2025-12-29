package docstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoc_GenerateDocId(t *testing.T) {
	set := make(map[string]struct{})

	for i := 0; i < 1000; i++ {
		id := GenerateDocId()

		_, exists := set[id.String()]
		assert.False(t, exists, "found duplicated generated doc id")

		set[id.String()] = struct{}{}
	}
}
