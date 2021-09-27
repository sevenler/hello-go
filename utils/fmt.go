package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(format string, args ...interface{}){
	argJsons := make([]interface{}, 0)
	for i := 0; i<len(args); i++{
		s, e := json.Marshal(args[i])
		fmt.Printf("--- %v ----\n", args[i])
		if e != nil{
			fmt.Printf("error format %v", e)
		}
		argJsons = append(argJsons, string(s))
	}
	fmt.Printf(format, argJsons...)
}