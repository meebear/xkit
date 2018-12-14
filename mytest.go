package main

import (
    "xkit/mytest/clip"
    //"xkit/mytest/packd"
    "fmt"
    "time"
    "net"
)

type config struct {
    port int
    port2 int
    verbose int
    enabled bool
    target string
    target2 string
    dura time.Duration
    ip net.IP
}
var cfg config
var cfg2 config

func main() {
    packd.WalkDir(".")
    return

    clip.ArgOption(&cfg.port, 'p', "port", "port-0")
    clip.ArgOption(&cfg.port2, 'P', "Port", "port-1")
    clip.IncrOption(&cfg.verbose, 'v',
            "verbose-and-long-option-name",
            "verbose level is a very long and unnecessarily long option with" +
            "an again very very long and hard to read description in English" +
            "written by a Chinese guy who's English is still in very preliminery" +
            "level...")
    clip.FlagOption(&cfg.enabled, 'e', "enable", "enable")
    clip.Positional(&cfg.target, "target", "target")
    clip.Positional(&cfg.target2, "target2", "target2")
    o := clip.ArgOption(&cfg.dura, 'd', "dura", "duration")
    o.Hide()
    clip.ArgOption(&cfg.ip, 'i', "ip", "duration")

    c := clip.SubCommand("cmd", "next level", nil)
    c.Hide()
    c.ArgOption(&cfg2.port, 'p', "port", "port-0")
    c.ArgOption(&cfg2.port2, 'P', "Port", "port-1")
    c.IncrOption(&cfg2.verbose, 'v', "verbose", "verbose level")
    c.FlagOption(&cfg2.enabled, 'h', "enable", "enable")
    c.Positional(&cfg2.target, "target", "target")
    c.Positional(&cfg2.target2, "target2", "target2")

    if _, err := clip.Parse(); err != nil {
        fmt.Println(">> ", err)
    }

    fmt.Printf("\n\ncfg:  %#v\n", cfg)
    fmt.Printf("cfg2: %#v\n", cfg2)
    fmt.Printf("done\n")
}
