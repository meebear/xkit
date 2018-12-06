package main

import (
    "xkit/mytest/clip"
    "fmt"
)

var num int

func main() {
    clip.ArgOption(&num, "num", "", "just a number")

    clip.PrintHelpCommand(&clip.RootCmd)

    clip.Parse()

    fmt.Println("done")
}
