package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/yangtudou/nb-action/actions/registry_sync"
)

func main() {

	if len(os.Args) < 2 {

		printUsage()

		os.Exit(1)
	}

	switch os.Args[1] {

	case "sync":

		runSync()

	default:

		fmt.Fprintf(
			os.Stderr,
			"unknown command: %s\n",
			os.Args[1],
		)

		printUsage()

		os.Exit(1)
	}
}

func runSync() {

	opt, err := registry_sync.ParseOptions(
		os.Args[2:],
	)

	if err != nil {

		fmt.Fprintf(
			os.Stderr,
			"invalid arguments: %v\n",
			err,
		)

		os.Exit(1)
	}

	if opt.DstPrefix == "" {

		fmt.Fprintln(
			os.Stderr,
			"missing required flag: --dst-prefix",
		)

		os.Exit(1)
	}

	result, err := registry_sync.Run(
		context.Background(),
		opt,
	)

	if err != nil {

		fmt.Fprintf(
			os.Stderr,
			"registry-sync failed: %v\n",
			err,
		)

		os.Exit(1)
	}

	printJSON(
		result,
	)
}

func printJSON(
	value interface{},
) {

	data, err := json.Marshal(
		value,
	)

	if err != nil {

		fmt.Fprintf(
			os.Stderr,
			"encode result failed: %v\n",
			err,
		)

		os.Exit(1)
	}

	fmt.Println(
		string(data),
	)
}

func printUsage() {

	fmt.Fprintln(
		os.Stderr,
		"Usage:",
	)

	fmt.Fprintln(
		os.Stderr,
		"  registry-sync sync [options]",
	)
}
