package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"net/http"
	"io"
)
const prefix = "https://kdb.tsukuba.ac.jp/syllabi/2023/%s/jpn"

func main() {
	app := &cli.App{
		Name: "kdbfetch",
		Usage: "fetch a specific subject",
		Action: fetch,
		UsageText: "kdbfetch [global options] [subject id]",
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func fetch(c *cli.Context) error {
	number := c.Args().Get(0)

	client := new(http.Client)
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf(prefix, number),
		nil,
	)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error: status: %d\n", resp.StatusCode)
	}
	io.Copy(os.Stdout, resp.Body)

	return nil
}
