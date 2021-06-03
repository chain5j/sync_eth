package main

import (
	"github.com/chain5j/sync_eth/cmd"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // cpu
	cmd.Execute()
}
