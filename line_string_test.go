package s2util

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

func TestLineString(t *testing.T) {

	const xmlstrings = `<ISO6709>+12+123/+12.3+123.4/+12.34+123.43/</ISO6709>`

	var ISO6709 LineString

	err := xml.Unmarshal([]byte(xmlstrings), &ISO6709)
	if err != nil {
		t.Error(err)
	}

	expct := NewPoint(NewAngle(12, 1), NewAngle(123, 1), nil)
	if ISO6709.MultiPoint[0] != expct {
		t.Errorf("Not match got %#v expct %#v", ISO6709.MultiPoint[0], expct)
	}

	expct = NewPoint(NewAngle(12.3, 0.1), NewAngle(123.4, 0.1), nil)
	if ISO6709.MultiPoint[1] != expct {
		t.Errorf("Not match got %#v expct %#v", ISO6709.MultiPoint[0], expct)
	}

	expct = NewPoint(NewAngle(12.34, 0.01), NewAngle(123.43, 0.01), nil)
	if ISO6709.MultiPoint[2] != expct {
		t.Errorf("Not match got %#v expct %#v", ISO6709.MultiPoint[0], expct)
	}

	b, err := json.Marshal(&ISO6709)
	if err != nil {
		t.Error(err)
	}

	expcts := `[[123,12],[123.4,12.3],[123.43,12.34]]`
	if string(b) != expcts {
		t.Errorf("Wrong got %s expct %s", string(b), expcts)
	}
}
