package main

import (
	"github.com/vietanhduong/ewol/cmd"
	"github.com/vietanhduong/ewol/pkg/cli"
)

func main() { cli.Execute(cmd.New()) }
