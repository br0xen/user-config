package main

import (
	"fmt"
	"os"

	userConfig "github.com/br0xen/user-config"
)

const AppName = "cfgedit"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	whichConfig := os.Args[1]
	cfg, err := userConfig.NewConfig(whichConfig)
	if err != nil {
		fmt.Println("Couldn't find config directory: " + whichConfig)
		os.Exit(1)
	}
	op := "list"
	if len(os.Args) >= 3 {
		op = os.Args[2]
	}
	switch op {
	case "list":
		fmt.Println(cfg.GetKeyList())
	}
}

func printHelp() {
	fmt.Println("Usage: " + AppName + " <which config> <operation>")
	fmt.Println("  <which-config> is ~/.config/<which-config>")
	fmt.Println("  <operation> can just be 'list' right now")
}
