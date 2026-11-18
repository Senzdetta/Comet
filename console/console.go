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
        "--version": Version{},
        "--help": Help{},
        "--uwu": Uwu{},
        "--portscan": Portscan{},
        "--pathfuzzer": Pathfuzzer{},
        "--userhunt": Userhunt{},
        "--paramscan": Paramscan{},
        "--tracelink": Tracelink{},
        "--header": Header{},
        "--subdomain": Subdomain{},
        "--hostsearch": Hostsearch{},
        "--dnslookup": Dnslookup{},
        "--sqlscan": Sqlscan{},
        "--subtakeover": Subtakeover{},
        "--corsscan": Corsscan{},
    }

    if cmd, ok := commands[args[1]]; ok {
        cmd.Execute(args)
    } else {
        invinput.InvalidOption(args[1])
        os.Exit(1)
    }
}

// Copyright (c) 2026 Senzdetta