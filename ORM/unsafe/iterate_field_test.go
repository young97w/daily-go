package unsafe

import "testing"

func TestIterate(t *testing.T) {
	testCases := []struct {
		name   string
		entity any
	}{
		{
			name:   "userV0",
			entity: userV0{},
		},
		{
			name:   "userV1",
			entity: userV1{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iterateFields(tc.entity)
		})
	}
}

type userV0 struct {
	Name    string
	Age     int32
	Alias   []string
	Address string
}

type userV1 struct {
	Name    string
	Age     int32
	AgeV1   int32
	Alias   []string
	Address string
}
