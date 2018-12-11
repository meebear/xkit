package main

import (
    "xkit/mytest/clip"
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
    clip.ArgOption(&cfg.port, 'p', "port", "port-0")
    clip.ArgOption(&cfg.port2, 'P', "Port", "port-1")
    clip.IncrOption(&cfg.verbose, 'v', "verbose", "verbose level")
    clip.FlagOption(&cfg.enabled, 'e', "enable", "enable")
    clip.Positional(&cfg.target, "target", "target")
    clip.Positional(&cfg.target2, "target2", "target2")
    clip.ArgOption(&cfg.dura, 'd', "dura", "duration")
    clip.ArgOption(&cfg.ip, 'i', "ip", "duration")

    c := clip.SubCommand("cmd", "next level")
    c.ArgOption(&cfg2.port, 'p', "port", "port-0")
    c.ArgOption(&cfg2.port2, 'P', "Port", "port-1")
    c.IncrOption(&cfg2.verbose, 'v', "verbose", "verbose level")
    c.FlagOption(&cfg2.enabled, 'e', "enable", "enable")
    c.Positional(&cfg2.target, "target", "target")
    c.Positional(&cfg2.target2, "target2", "target2")

    clip.PrintHelpCommand(&clip.RootCmd)
    fmt.Printf("\n\n")

    if _, err := clip.Parse(); err != nil {
        fmt.Println(">> ", err)
    }

    fmt.Printf("\n\ncfg:  %#v\n", cfg)
    fmt.Printf("cfg2: %#v\n", cfg2)
    fmt.Printf("done\n")
}
