// Package bits implements the encoding and decoding of bits
package bits

import (
    "io"
    "errors"
    "bytes"
)

func bitToInt(bit bool) int {
    n := 0
    if bit {
        n = 1
    }
    return n
}

type Encoder struct {
    bits []bool
    bytes []byte
}

// NewEncoder creates a new encoder
func NewEncoder() Encoder {
    return Encoder{make([]bool, 0), make([]byte, 0)}
}

// PutBit will store the bit.
func (enc *Encoder) PutBit(bit bool) {
    enc.bits = append(enc.bits, bit)
    if len(enc.bits) == 8 {
        var b byte
        for i := 7; i >= 0; i-- {
            n := bitToInt(enc.bits[7 - i])
            b += byte(n << uint8(i))
        }
        enc.bytes = append(enc.bytes, b)
        enc.bits = make([]bool, 0)
    }
}

// PutBits will store all the bits in the slice given.
func (enc *Encoder) PutBits(bits []bool) {
    for _, bit := range bits {
        enc.PutBit(bit)
    }
}

// RemainderBits returns the number of individual bits that are being store. If there has been a multiple of 8 bits stored 
// then it will return 0, otherwise it will return the remainder of the number of bits stored divided by 8. For example if 9 bits
// are stored then it will return 1.
//
// This may be useful to check whether padding is needed, without having to handle the error from encoder.Write
//
// If you just want the number of bits stored, use encoder.Len
func (enc Encoder) RemainderBits() int {
    return len(enc.bits)
}

// Len returns the number of bits stored
func (enc Encoder) Len() int {
    return len(enc.bytes) * 8 + len(enc.bits)
}

func (enc *Encoder) pad(padBit bool) bool {
    length := len(enc.bits)
    if length == 0 {
        return false
    }
    for bitsNeeded := 8 - length; bitsNeeded > 0; bitsNeeded-- {
        enc.PutBit(padBit)
    }
    return true
}

// Write will write the stored bits to writer. If there is not a multiple of 8 bits to store, ie. there is not a round number 
// of bytes to write, then it will return an error. Otherwise it will return the error given by writing to writer.
func (enc *Encoder) Write(writer io.Writer) error {
    length := len(enc.bits)
    if length != 0 {
        return errors.New("There is not a whole numbers of bytes available to be written (maybe set pad to true?)")
    }
    _, err := writer.Write(enc.bytes)
    return err
}

// WriteAndPad will write the stored bits to writer. If there is not a round number of bytes to store then it will add
// the required number of bits, with value padBit, to the end. It will return whether or not it had to pad and the error 
// given by writing to writer.
func (enc *Encoder) WriteAndPad(writer io.Writer, padBit bool) (bool, error) {
    padded := enc.pad(padBit)
    err := enc.Write(writer)
    return padded, err
}

type Decoder struct {
    reader io.Reader
    data []bool
}

// NewDecoder creates a new decoder.
func NewDecoder(reader io.Reader) Decoder {
    return Decoder{reader, make([]bool, 0)}
}

// NewDecoderFromBytes creates a new decoder which reads from the slice b
func NewDecoderFromBytes(b []byte) Decoder {
    return Decoder{bytes.NewBuffer(b), make([]bool, 0)}
}

// GetBit will return one bit from the reader. It will return any errors from getting the data from the reader. It should be noted 
// that it may have consumed more data than expected from the reader, as readers only allow bytes to be read from them.
func (dec *Decoder) GetBit() (bool, error) {
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

// GetBits will get all of the bits from the reader and return any errors it encounters from accessing the reader
func (dec *Decoder) GetBits() ([]bool, error) {
    var bit bool
    var err error
    bits := make([]bool, 0)
    for err == nil {
        bit, err = dec.GetBit()
        bits = append(bits, bit)
    }
    if err == io.EOF {
        return bits[:len(bits) - 1], nil
    }
    return nil, err
}
