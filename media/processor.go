package media

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/OdaDaisuke/media-concatman/asset"
)

type Processor struct {
	framerate string
	setting   *asset.AssetSettings
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
	resourcePaths, err := p.initResources()
	if err != nil {
		return err
	}
	dist, err := p.concatResources(resourcePaths)
	if err != nil {
		return err
	}
	fmt.Printf("Output file -> %s", dist)
	return nil
}

func (p *Processor) saveMp4Files(filename string, imgOnlyDistPaths []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := 0; i < len(imgOnlyDistPaths); i++ {
		line := fmt.Sprintf("file %s\n", imgOnlyDistPaths[i])
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) concatResources(resourcePaths []string) (string, error) {
	distPath := "./dist.mp4"
	sourceFilename := "resources.txt"
	err := p.saveMp4Files(sourceFilename, resourcePaths)
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
	for i := 0; i < len(resourcePaths); i++ {
		os.Remove(resourcePaths[i])
		os.Remove(sourceFilename)
	}
	return distPath, nil
}

func (p *Processor) initResources() ([]string, error) {
	var distPaths []string
	for i := 0; i < len(p.setting.Resources); i++ {
		resource := p.setting.Resources[i]
		onlyImgVideoPath, err := p.initWithImage(i, resource)
		if err != nil {
			return nil, err
		}
		videoPaths, err := p.concatAudio(i, onlyImgVideoPath, resource)
		if err != nil {
			return nil, err
		}
		for _, videoPath := range videoPaths {
			distPaths = append(distPaths, videoPath)
		}
	}
	return distPaths, nil
}

func (p *Processor) initWithImage(i int, resource asset.Resource) (string, error) {
	inputFilePath := p.setting.Context.ImageAssetPath + "/" + resource.ImageFile.Filename
	ffmpegCmd, err := NewFFMpeg(inputFilePath, nil)
	if err != nil {
		return "", err
	}
	dist := "dist_img_" + strconv.Itoa(i) + ".mp4"

	ffmpegCmd.SetArgs("-loop", "1")
	ffmpegCmd.AllowOverwrite()
	ffmpegCmd.SetArgs("-vcodec", "libx264")
	ffmpegCmd.SetArgs("-r", p.framerate)
	ffmpegCmd.SetArgs("-t", "100")
	ffmpegCmd.SetArgs("-s", "1280x720")
	ffmpegCmd.SetArgs(dist)
	_, err = ffmpegCmd.Execute()
	return dist, err
}

func (p *Processor) concatAudio(resourceIndex int, srcVideoPath string, resource asset.Resource) ([]string, error) {
	var distPaths []string
	var audioFiles []string
	if resource.AudioFile.IsDir == true {
		files, err := ioutil.ReadDir(p.setting.Context.AudioAssetPath + resource.AudioFile.Dirname)
		if err != nil {
			return nil, err
		}
		for i, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.Contains(file.Name(), ".wav") || strings.Contains(file.Name(), ".mp3") {
				distPaths = append(distPaths, fmt.Sprintf("dist_resource_%d_%d.mp4", resourceIndex, i))
				audioFiles = append(audioFiles, file.Name())
			}
		}
	} else {
		distPaths = append(distPaths, fmt.Sprintf("dist_resource_%d.mp4", resourceIndex))
		audioFiles = append(audioFiles, resource.AudioFile.Filename)
	}

	for i, distPath := range distPaths {
		ffmpegCmd, err := NewFFMpeg(srcVideoPath, nil)
		if err != nil {
			return nil, err
		}
		filepath := p.setting.Context.AudioAssetPath + "/" + audioFiles[i]
		ffmpegCmd.SetArgs("-i", filepath)
		ffmpegCmd.SetArgs("-c:v", "copy")
		ffmpegCmd.SetArgs("-c:a", "aac")
		ffmpegCmd.SetArgs("-map", "0:v:0")
		ffmpegCmd.SetArgs("-map", "1:a:0")
		ffmpegCmd.AllowOverwrite()
		ffmpegCmd.SetArgs(distPath)
		_, err = ffmpegCmd.Execute()
		if err != nil {
			return nil, err
		}
	}
	os.Remove(srcVideoPath)
	return distPaths, nil
}

func (p *Processor) reductionNoise() {
	// ffmpeg -i <input_file> -af "highpass=f=200, lowpass=f=3000" <output_file>
}
