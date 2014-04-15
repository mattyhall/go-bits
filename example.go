package main

import (
    "github.com/mattyhall/go-bits/bits"
    "bytes"
    "fmt"
)

func main() {
    buff := bytes.NewBuffer(make([]byte, 0))

    encoder := bits.NewEncoder()
    // Store the 7 bits 1010111
    encoder.PutBits([]bool{false, true, false, true, true, true, true})
    encoder.PutByte(211)
    fmt.Println(encoder.Len(), encoder.RemainderBits())
    // Write the bits to the buffer. As there is only 7 bits then one extra one will be added
    // It will be a 0 as false is passed to WriteAndPad
    encoder.WriteAndPad(buff, false)

    // Could use bits.NewDecoderFromBytes(buff.Bytes())
    decoder := bits.NewDecoder(buff)
    // Read the bits back
    bits, err := decoder.GetBits(7)
    fmt.Println(bits, err)
    var b byte
    // Read the byte back
    b, err = decoder.GetByte()
    fmt.Println(b, err)
    // We had to pad so there should be one extra 'false' bit
    var rest []bool
    rest, err = decoder.GetBits(10)
    fmt.Println(rest, err)
}
