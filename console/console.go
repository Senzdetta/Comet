// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "github.com/Senzdetta/Comet/utils/invinput"
)

func Console(input string) {
    args := os.Args
    if len(args) < 2 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    commands := map[string]Command{
        "--version": Version{},             // ok
        "--help": Help{},                   // ok
        "--uwu": Uwu{},                     // ok
        "--portscan": Portscan{},           // ok
        "--pathfuzzer": Pathfuzzer{},       // ok
        "--userhunt": Userhunt{},           // ok
        "--paramscan": Paramscan{},         // ok
        "--tracelink": Tracelink{},
        "--header": Header{},               // ok
        "--subdomain": Subdomain{},         // ok
        "--hostsearch": Hostsearch{},       // ok
        "--dnslookup": Dnslookup{},         // ok
        "--sqlscan": Sqlscan{},             // ok
        "--subtakeover": Subtakeover{},     // ok
        "--corsscan": Corsscan{},           // ok
    }

    if cmd, ok := commands[args[1]]; ok {
        cmd.Execute(args)
    } else {
        invinput.InvalidOption(args[1])
        os.Exit(1)
    }
}

// Copyright (c) 2026 Senzdetta