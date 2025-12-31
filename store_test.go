package docstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	Clear()

	jimId, jim := AddTestDoc(t, "jim", 22)
	AssertDoc(t, jimId, jim)

	joeId, joe := AddTestDoc(t, "joe", 32)
	AssertDoc(t, joeId, joe)

	bobId, bob := AddTestDoc(t, "bob", 42)
	AssertDoc(t, bobId, bob)

	err := Delete(jimId)
	require.NoError(t, err)
	AssertNoDoc(t, jimId)
	AssertDoc(t, joeId, joe)
	AssertDoc(t, bobId, bob)

	jimId, jim = AddTestDoc(t, "jim", 52)
	AssertDoc(t, jimId, jim)
	AssertDoc(t, joeId, joe)
	AssertDoc(t, bobId, bob)

	err = Put(GenerateDocId(), 47)

	count := 0

	for _, _ = range AllDocuments() {
		count++
	}

	assert.Equal(t, 4, count)

	count = 0

	for _, _ = range AllDocumentsOf[TestDoc]() {
		count++
	}

	assert.Equal(t, 3, count)

	Clear()

	AssertNoDoc(t, jimId)
	AssertNoDoc(t, joeId)
	AssertNoDoc(t, bobId)
}
