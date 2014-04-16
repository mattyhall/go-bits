package bits

import (
	"bytes"
	"testing"
)

func myError(t *testing.T, txt string, expected interface{}, got interface{}) {
	t.Errorf(txt+": Expected %v, got %v", expected, got)
}

func testWriteLength(t *testing.T, writer string, buf *bytes.Buffer, expected int) {
	l := len(buf.Bytes())
	if l != expected {
		myError(t, "encoder."+writer+"() should write the correct number of bytes", expected, l)
	}
}

func TestPutBits(t *testing.T) {
	tests := [][]bool{[]bool{}, []bool{true}, []bool{true, true, true, true, true, true, true, true},
		[]bool{true, true, true, true, true, true, true, true, true}}
	for _, bs := range tests {
		encoder := NewEncoder()
		encoder.PutBits(bs...)
		if len(bs) != encoder.Len() {
			myError(t, "encoder.Len() should equal the number of bits written", len(bs), encoder.Len())
		}
		if len(bs)%8 != encoder.RemainderBits() {
			myError(t, "encoder.RemainderBits() should equal bits % 8", len(bs)%8, encoder.RemainderBits())
		}

		if len(bs)%8 == 0 {
			buf := bytes.NewBuffer(make([]byte, 0))
			err := encoder.Write(buf)
			if err != nil {
				myError(t, "encoder.Write() should not give an error if there is a round number of bytes", nil, err)
			}

			var padded bool
			buf = bytes.NewBuffer(make([]byte, 0))
			padded, err = encoder.WriteAndPad(buf, false)
			if padded {
				t.Error("encoder.WriteAndPad should not pad if there is a round number of bits")
			}
			testWriteLength(t, "WriteAndPad", buf, len(bs)/8)
		} else {
			buf := bytes.NewBuffer(make([]byte, 0))
			err := encoder.Write(buf)
			if err == nil {
				t.Error("encoder.Write() should give an error if there isn't a round number of bytes")
			}

			buf = bytes.NewBuffer(make([]byte, 0))
			var padded bool
			padded, err = encoder.WriteAndPad(buf, false)
			if !padded {
				t.Error("encoder.WriteAndPad should pad for an odd number of bits")
			}
			if err != nil {
				myError(t, "encoder.WriteAndPad should not throw an error", nil, err)
			}
		}
	}
}
