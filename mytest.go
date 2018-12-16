package main

import (
    "xkit/mytest/clip"
    "fmt"
)

func main() {
    clip.ProgDescription("mytest v0.0")
    packCmdInit(&clip.RootCmd)

    if _, err := clip.Parse(); err != nil {
        fmt.Println(">> ", err)
    }

    fmt.Printf("done\n")
}
