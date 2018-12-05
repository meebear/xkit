package main

import (
    "xkit/mytest/clip"
    "fmt"
)

var num int

func main() {
    o := clip.ArgOption(&num, "num", "just a number")

    fmt.Printf("%v  -- %#v\n", &num, o)
}
