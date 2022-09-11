package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	displayLogo()

	app := &cli.App{
		Name:  "codegen",
		Usage: "generate the code",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func displayLogo() bool {
	if len(os.Args) == 1 {
		color.Cyan(logo)
		return true
	} else if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "h" || os.Args[1] == "help" {
			color.Cyan(logo)
			return true
		}
	}
	return false
}
