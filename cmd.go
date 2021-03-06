package main

import (
	"github.com/OdaDaisuke/media-concatman/asset"
	"github.com/OdaDaisuke/media-concatman/media"
	"github.com/goccy/go-yaml"
)

func main() {
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
