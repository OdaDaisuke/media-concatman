package main

import (
	"os"

	"github.com/OdaDaisuke/media-concatman/asset"
	"github.com/OdaDaisuke/media-concatman/media"
	"github.com/goccy/go-yaml"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "media-concat"
	app.Usage = "media-concat"
	app.Version = "1.0.0"
	app.Action = func(c *cli.Context) {
		println("Number of arguments is", len(c.Args()))
		yamlBytes, err := asset.LoadFile()
		if err != nil {
			panic(err)
		}
		var v asset.AssetSettings
		if err := yaml.Unmarshal(yamlBytes, &v); err != nil {
			panic(err)
		}
		p := media.NewProcessor(v)
		if err = p.Run(); err != nil {
			panic(err)
		}
	}
	app.Run(os.Args)
}
