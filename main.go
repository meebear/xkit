package main

import (
    "xkit/mytest/clip"
    "fmt"
)

func main() {
    clip.ProgDescription("fazutil v0.1.0")
    awsCmdInit(&clip.RootCmd)

    var err error
    c, err := clip.Parse()
    fmt.Printf("%v, %v\n", c==nil, err==nil)
    if err == nil {
        if c.Run != nil {
            err = c.Run(c)
        } else {
            clip.HelpCommand(c, false)
        }
    }

    if err != nil {
        fmt.Printf("error: %s\n", err)
    }
}
