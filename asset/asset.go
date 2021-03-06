package asset

import "io/ioutil"

type Output struct {
	Filename string `yaml:"Filename"`
	Codec    string `yaml:"Codec"`
}

type Context struct {
	ImageAssetPath string `yaml:"ImageAssetPath"`
	AudioAssetPath string `yaml:"AudioAssetPath"`
	Duration       uint32 `yaml:"Duration"`
	Output         Output `yaml:"Output"`
}

type Resource struct {
	ImageFilename string `yaml:"ImageFilename"`
	AudioFilename string `yaml:"AudioFilename"`
}

type AssetSettings struct {
	Context   Context    `yaml:"Context"`
	Resources []Resource `yaml:"Resources"`
}

func LoadFile() ([]byte, error) {
	fname := "asset.yaml"
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		var b []byte
		return b, err
	}
	return bytes, nil
}
