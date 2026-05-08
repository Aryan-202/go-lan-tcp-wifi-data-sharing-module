package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
        Name:  "send",
        Usage: "this command will send file",
        Action: func(context.Context, *cli.Command) error {
            fmt.Println("send file")
            return nil
        },
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}