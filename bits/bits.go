// Package bits implements the encoding and decoding of bits
package bits

import (
    "io"
)

func bitToInt(bit bool) int {
    n := 0
    if bit {
        n = 1
    }
    return n
}

type Encoder struct {
    writer io.Writer
    data []bool
}

// NewEncoder creates a new encoder which will write its data to writer 
func NewEncoder(writer io.Writer) Encoder {
    return Encoder{writer, make([]bool, 0)}
}

// EncodeBit will encode a bit
func (enc *Encoder) EncodeBit(bit bool) bool {
    enc.data = append(enc.data, bit)
    if len(enc.data) == 8 {
        var b byte
        for i := 7; i >= 0; i-- {
            n := bitToInt(enc.data[7 - i])
            b += byte(n << uint8(i))
        }
        bytes := []byte{b}
        enc.writer.Write(bytes)
        enc.data = make([]bool, 0)
        return true
    }
    return false
}

func (enc *Encoder) EncodeBits(bits []bool) bool {
    allWritten := true
    for _, bit := range bits {
        allWritten = enc.EncodeBit(bit)
    }
    return allWritten
}

func (enc *Encoder) Flush(defaultBit bool) {
    length := len(enc.data)
    if length == 0 {
        return
    }
    for bitsNeeded := 8 - length; bitsNeeded > 0; bitsNeeded-- {
        enc.EncodeBit(defaultBit)
    }
}

type Decoder struct {
    reader io.Reader
    data []bool
}

func NewDecoder(reader io.Reader) Decoder {
    return Decoder{reader, make([]bool, 0)}
}

func (dec *Decoder) DecodeBit() (bool, error) {
    if len(dec.data) == 0 {
        bytes := make([]byte, 1)
        if _, err := dec.reader.Read(bytes); err != nil {
            return false, err
        }
        for i := 7; i >= 0; i-- {
            asBool := (bytes[0] >> uint8(i)) & 1 == 1
            dec.data = append(dec.data, asBool)
        }
    }
    b := dec.data[0]
    dec.data = dec.data[1:]
    return b, nil
}

func (dec *Decoder) DecodeBits() ([]bool, error) {
    var bit bool
    var err error
    bits := make([]bool, 0)
    for err == nil {
        bit, err = dec.DecodeBit()
        bits = append(bits, bit)
    }
    if err == io.EOF {
        return bits[:len(bits) - 1], nil
    }
    return nil, err
}
