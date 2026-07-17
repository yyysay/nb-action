package registry_sync

import (
	"flag"
	"io"
)

type Options struct {
	Src         string
	Base        string
	SrcPrefix   string
	DstPrefix   string
	SrcFlatten  bool
	DstFlatten  bool
	Platform    string
	Concurrency int
	Retries     int
	DryRun      bool
}

func ParseOptions(args []string) (*Options, error) {

	fs := flag.NewFlagSet(
		"registry-sync",
		flag.ContinueOnError,
	)

	fs.SetOutput(io.Discard)

	opt := &Options{}

	fs.StringVar(&opt.Src, "src", "", "source image")
	fs.StringVar(&opt.Base, "base", "images.txt", "image list")
	fs.StringVar(&opt.SrcPrefix, "src-prefix", "", "source prefix")
	fs.StringVar(&opt.DstPrefix, "dst-prefix", "", "destination prefix")

	fs.BoolVar(&opt.SrcFlatten, "src-flatten", false, "flatten source")
	fs.BoolVar(&opt.DstFlatten, "dst-flatten", false, "flatten target")

	fs.StringVar(&opt.Platform, "platform", "", "platform")

	fs.IntVar(&opt.Concurrency, "concurrency", 4, "workers")
	fs.IntVar(&opt.Retries, "retries", 3, "retry")

	fs.BoolVar(&opt.DryRun, "dry-run", false, "dry run")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return opt, nil
}
