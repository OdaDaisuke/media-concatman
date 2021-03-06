package media

import (
	"fmt"
	"os/exec"
)

type FFMpeg struct {
	*exec.Cmd
}

func NewFFMpeg(inputFilePath string, format *string) (*FFMpeg, error) {
	ffmpegCmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}
	if format != nil {
		return &FFMpeg{
			exec.Command(ffmpegCmdPath, "-f", *format, "-i", inputFilePath),
		}, nil
	}
	return &FFMpeg{
		exec.Command(ffmpegCmdPath, "-i", inputFilePath),
	}, nil
}

func (f *FFMpeg) AllowOverwrite() {
	f.SetArgs("-y")
}

func (f *FFMpeg) SetArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *FFMpeg) Execute() ([]byte, error) {
	fmt.Println(f.Args)
	return f.CombinedOutput()
}
