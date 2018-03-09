package main

import (
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	murmur "github.com/mono0x/murmur/lib"
	"golang.org/x/sync/errgroup"
)

//go:generate retool do go-assets-builder -o assets.go config

func handler() error {
	var eg errgroup.Group
	for _, file := range Assets.Files {
		if file.IsDir() {
			continue
		}

		path := file.Path
		if !(strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
			continue
		}

		eg.Go(func() error {
			file, err := Assets.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			config, err := murmur.LoadConfig(file)
			if err != nil {
				return err
			}

			source, err := config.NewSource()
			if err != nil {
				return err
			}

			sink, err := config.NewSink()
			if err != nil {
				return err
			}
			defer sink.Close()

			notifier := murmur.NewNotifier(source, sink)
			return notifier.Notify()
		})
	}
	return eg.Wait()
}

func main() {
	lambda.Start(handler)
}
