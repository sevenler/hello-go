package main

import (
	"encoding/json"
	"fmt"
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

func PrintJson(format string, args ...interface{}){
	argJsons := make([]interface{}, 0)
	for i := 0; i<len(args); i++{
		s, e := json.Marshal(args[i])
		if e != nil{
			fmt.Printf("error format %v", e)
		}
		argJsons = append(argJsons, string(s))
	}
	fmt.Printf(format, argJsons...)
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "rm <path>":
		PrintJson("Get args: %v\n", CLI)
	case "ls":
		PrintJson("Get args: %v\n", CLI)
	default:
		panic(ctx.Command())
	}
}