package docstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArchive(t *testing.T) {
	RegisterType[TestDoc]()

	Clear()

	jimId, jim := AddTestDoc(t, "jim", 22)
	joeId, joe := AddTestDoc(t, "joe", 32)
	bobId, bob := AddTestDoc(t, "bob", 42)

	filename := TempFilename()

	err := WriteAllToFile(filename)
	require.NoError(t, err)

	Clear()
	AssertNoDoc(t, jimId)
	AssertNoDoc(t, joeId)
	AssertNoDoc(t, bobId)

	err = ReadAllFromFile(filename)
	require.NoError(t, err)

	AssertDoc(t, jimId, jim)
	AssertDoc(t, joeId, joe)
	AssertDoc(t, bobId, bob)
}
