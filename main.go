package main

import (
    "encode/bits"
    "bytes"
    "fmt"
)

func main() {
    buff := bytes.NewBuffer(make([]byte, 0))
    encoder := bits.NewEncoder()
    encoder.PutBits([]bool{true, false, true, false, true, true, true})
    encoder.WriteAndPad(buff, false)
    decoder := bits.NewDecoder(buff)
    bs, err := decoder.GetBits()
    fmt.Println(bs, err, buff)
}
