package docstore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	RegisterType[TestDoc]()

	Clear()

	jimId, jim := AddTestDoc(t, "jim", 22)
	joeId, joe := AddTestDoc(t, "joe", 32)
	bobId, bob := AddTestDoc(t, "bob", 42)

	var buffer bytes.Buffer

	err := WriteAll(&buffer)
	require.NoError(t, err)

	Clear()
	AssertNoDoc(t, jimId)
	AssertNoDoc(t, joeId)
	AssertNoDoc(t, bobId)

	err = ReadAll(&buffer)
	require.NoError(t, err)

	AssertDoc(t, jimId, jim)
	AssertDoc(t, joeId, joe)
	AssertDoc(t, bobId, bob)
}
