package main

import (
    "encode/bits"
    "bytes"
    "fmt"
)

func main() {
    buff := bytes.NewBuffer(make([]byte, 0))
    encoder := bits.NewEncoder(buff)
    encoder.EncodeBits([]bool{true, false, true, false, true, true, true, })
    encoder.Flush(false)
    decoder := bits.NewDecoder(buff)
    bs, err := decoder.DecodeBits()
    fmt.Println(bs, err, buff)
}
