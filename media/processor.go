package media

import (
	"fmt"
	"os"
	"strconv"

	"github.com/OdaDaisuke/media-concatman/asset"
)

type Processor struct {
	framerate         string
	setting           *asset.AssetSettings
	distImgOnlyFiles  []string
	distResourceFiles []string
}

func NewProcessor(setting asset.AssetSettings) *Processor {
	return &Processor{
		framerate: "1",
		setting:   &setting,
	}
}

func (p *Processor) Run() error {
	if p.setting == nil {
		panic("no setting")
	}
	if len(p.setting.Resources) == 0 {
		return nil
	}
	err := p.initResources()
	if err != nil {
		return err
	}
	dist, err := p.concatResources()
	if err != nil {
		return err
	}
	fmt.Printf("Output file -> %s", dist)
	return nil
}

func (p *Processor) saveMp4Files(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := 0; i < len(p.setting.Resources); i++ {
		line := fmt.Sprintf("file %s\n", p.distResourceFiles[i])
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) concatResources() (string, error) {
	distPath := "./dist.mp4"
	sourceFilename := "resources.txt"
	err := p.saveMp4Files(sourceFilename)
	if err != nil {
		return "", err
	}

	// ffmpeg -f concat -i {filename}.txt -c copy dist.mp4
	format := "concat"
	ffmpegCmd, err := NewFFMpeg(sourceFilename, &format)
	if err != nil {
		return "", err
	}
	ffmpegCmd.AllowOverwrite()
	ffmpegCmd.SetArgs("-c", "copy")
	ffmpegCmd.SetArgs(distPath)
	out, err := ffmpegCmd.Execute()
	fmt.Println(string(out))
	if err != nil {
		return "", err
	}
	for i := 0; i < len(p.setting.Resources); i++ {
		os.Remove(p.distResourceFiles[i])
		os.Remove(p.distImgOnlyFiles[i])
		os.Remove(sourceFilename)
	}
	return distPath, nil
}

func (p *Processor) initResources() error {
	for i := 0; i < len(p.setting.Resources); i++ {
		resource := p.setting.Resources[i]
		err := p.initWithImage(i, resource)
		if err != nil {
			return err
		}
		err = p.concatAudio(i, resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) initWithImage(i int, resource asset.Resource) error {
	inputFilePath := p.setting.Context.ImageAssetPath + "/" + resource.ImageFilename
	ffmpegCmd, err := NewFFMpeg(inputFilePath, nil)
	if err != nil {
		return err
	}
	dist := "dist_img_" + strconv.Itoa(i) + ".mp4"
	p.distImgOnlyFiles = append(p.distImgOnlyFiles, dist)

	ffmpegCmd.SetArgs("-loop", "1")
	ffmpegCmd.AllowOverwrite()
	ffmpegCmd.SetArgs("-vcodec", "libx264")
	ffmpegCmd.SetArgs("-r", p.framerate)
	ffmpegCmd.SetArgs("-t", "100")
	ffmpegCmd.SetArgs("-s", "1280x720")
	ffmpegCmd.SetArgs(dist)
	_, err = ffmpegCmd.Execute()
	return err
}

func (p *Processor) concatAudio(i int, resource asset.Resource) error {
	ffmpegCmd, err := NewFFMpeg(p.distImgOnlyFiles[i], nil)
	if err != nil {
		return err
	}
	dist := "dist_resource_" + strconv.Itoa(i) + ".mp4"
	p.distResourceFiles = append(p.distResourceFiles, dist)
	filepath := p.setting.Context.AudioAssetPath + "/" + resource.AudioFilename
	ffmpegCmd.SetArgs("-i", filepath)
	ffmpegCmd.SetArgs("-c:v", "copy")
	ffmpegCmd.SetArgs("-c:a", "aac")
	ffmpegCmd.SetArgs("-map", "0:v:0")
	ffmpegCmd.SetArgs("-map", "1:a:0")
	ffmpegCmd.AllowOverwrite()
	ffmpegCmd.SetArgs(dist)
	_, err = ffmpegCmd.Execute()
	return err
}

func (p *Processor) reductionNoise() {
	// ffmpeg -i <input_file> -af "highpass=f=200, lowpass=f=3000" <output_file>
}
