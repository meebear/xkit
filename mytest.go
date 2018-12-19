package main

import (
    "xkit/mytest/clip"
    "fmt"
)

func main() {
    clip.ProgDescription("mytest v0.0")
    packCmdInit(&clip.RootCmd)
    awsCmdInit(&clip.RootCmd)

    var err error
    c, err := clip.Parse()
    if err == nil {
        err = c.Run()
    }

    if err == nil {
        fmt.Printf("done\n")
    } else {
        fmt.Printf("error: %s\n", err)
    }
}
