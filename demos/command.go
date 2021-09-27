package main

import (
	"git.in.zhihu.com/zhihu/hello/utils"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Rm struct {
		Force bool   `help:"Force removal."`
		Value string `help:"Recursively remove files."`

		Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
	} `cmd:"" help:"Remove files."`

	Ls struct {
		Paths []string `arg:"" optional:"" name:"path" help:"Paths to list." type:"path"`
	} `cmd:"" help:"List paths."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "rm <path>":
		utils.PrintJson("Get args: %v\n", CLI)
	case "ls":
		utils.PrintJson("Get args: %v\n", CLI)
	default:
		panic(ctx.Command())
	}
}