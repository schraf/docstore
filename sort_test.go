package docstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSort_ByField(t *testing.T) {
	docs := []*Document[TestDoc]{
		{Id: GenerateDocId(), Data: TestDoc{Name: "one", Age: 1}},
		{Id: GenerateDocId(), Data: TestDoc{Name: "three", Age: 3}},
		{Id: GenerateDocId(), Data: TestDoc{Name: "two", Age: 2}},
	}

	sorted := sortDocuments(
		docs,
		func(a, b *Document[TestDoc]) bool {
			return a.Data.Age < b.Data.Age
		},
	)

	assert.Equal(t, 1, sorted[0].Data.Age)
	assert.Equal(t, 2, sorted[1].Data.Age)
	assert.Equal(t, 3, sorted[2].Data.Age)
}

func TestSort_WithTimestamps(t *testing.T) {
	now := time.Now()
	docs := []*Document[TestDoc]{
		{Id: NewDocId("c"), CreatedAt: now.Add(2 * time.Second)},
		{Id: NewDocId("a"), CreatedAt: now},
		{Id: NewDocId("b"), CreatedAt: now.Add(1 * time.Second)},
	}

	sorted := sortDocuments(
		docs,
		func(a, b *Document[TestDoc]) bool {
			return a.CreatedAt.Before(b.CreatedAt)
		},
	)

	assert.Equal(t, "a", sorted[0].Id.String())
	assert.Equal(t, "b", sorted[1].Id.String())
	assert.Equal(t, "c", sorted[2].Id.String())
}
