package s2util

import (
	"testing"
)

func TestMultiPoint_UnmarshalText(t *testing.T) {
	var mp MultiPoint

	err := mp.UnmarshalText([]byte(`+352139+1384339+3776/`))

	if err != nil {
		t.Errorf("MultiPoint error %#v", err)
	}
}
