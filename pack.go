package main

import (
    "xkit/mytest/clip"
    //"xkit/mytest/packd"
    "fmt"
)

var path string

func packCmdInit(c *clip.Command) {
    sc := c.SubCommand("pack", "pack/unpack a directory/file", packRun).Hide()
    sc.ArgOption(&path, 'p', "path", "path to pack/unpack").MustSet()
}

func packRun() error {
    fmt.Printf("packRun running... path=%s\n", path)
    return nil
}
