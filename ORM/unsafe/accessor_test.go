package unsafe

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestField(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	u := &user{
		Name: "young",
		Age:  18,
	}

	a := NewUnsafeAccessor(u)
	res, err := a.field("Age")
	require.NoError(t, err)
	require.Equal(t, 18, res)

	err = a.setField("Age", 20)
	require.NoError(t, err)
	require.Equal(t, 20, u.Age)
}
