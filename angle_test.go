package s2util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAngle(t *testing.T) {
	var a Angle
	b := []byte(`270`)
	err := json.Unmarshal(b, &a)
	if err != nil {
		t.Error(err.Error())
	}

	expct := `270`
	if a.String() != expct {
		t.Errorf("Got %s expect %s", a.String(), expct)
	}

	a = NewAngle(270, 1)

	bb, err := json.Marshal(&a)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf(string(bb))
	// if !bytes.Equal(b, bb) {
	// 	t.Errorf(`Got "%s" expect "%s"`, string(b), string(bb))
	// }
}
