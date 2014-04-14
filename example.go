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
    encoder.PutBits([]bool{true, false, true, false, true, true, true})
    // Write the bits to the buffer. As there is only 7 bits then one extra one will be added
    // It will be a 0 as false is passed to WriteAndPad
    encoder.WriteAndPad(buff, false)

    decoder := bits.NewDecoder(buff)
    // Read the bits back, should get 8: 10101110
    bs, err := decoder.GetBits()
    fmt.Println(bs, err, buff)
}
