package converter

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	tplDIR_FAKE = "files"
	tplDIR_REAL = "output"

	pthFFMPEG = "/usr/local/opt/ffmpeg/bin" // I have FFmpeg installed
	binFFMPEG = "ffmpeg"                    // with homebrew as `brew install ffmpeg`
)

type conv struct {
	Done chan bool
	Err  error
	Info []string
	pth  string
}

func New(pth string) *conv {
	c := conv{}
	c.Done = make(chan bool)
	c.pth = pth
	return &c
}

func (c *conv) Check() bool {
	var (
		fin      = c.pth
		setError = func() bool {
			c.Err = fmt.Errorf("file path error, for '%s', at %v", c.pth, time.Now())
			return false
		}
	)

	if fin == "" {
		return setError()
	}

	lst := strings.Split(fin, "/")
	if len(lst) < 2 {
		return setError()
	}
	val := strings.TrimSpace(lst[0])
	if val == "." || val == "" { // remove first element if it '.' or empty
		lst = lst[1:]
	}
	lst[0] = tplDIR_REAL // replace user entry for static files with real dir name

	fl, err := os.Executable()
	if err != nil {
		log.Println(err)
		return setError()
	}
	lst = append([]string{filepath.Dir(fl)}, lst...) // prepend exec. path
	fin = filepath.Join(lst...)                      // make file path as string

	if _, err := os.Stat(fin); errors.Is(err, os.ErrNotExist) {
		return setError()
	}

	c.pth = fin
	return true
}

func (c *conv) Run() string {
	runner := func(msg string) {
		var (
			fin     = c.pth
			pth, fl = filepath.Split(fin)
		)
		c.Info = append(c.Info, msg)

		// replace current extention by '.mp3'
		out := filepath.Join(pth, (fl[0:len(fl)-len(filepath.Ext(fl))] + ".mp3"))

		cmd := exec.Command(filepath.Join(pthFFMPEG, binFFMPEG), []string{
			"-hide_banner",
			"-i", fin,
			"-y", out}...)
		cmd.Dir = pth
		output, _ := cmd.CombinedOutput()

		// actually run FFmpeg for convert media to MP3 file
		if err := cmd.Run(); err != nil {
			c.Info = append(c.Info, err.Error())
		}
		c.Info = append(c.Info, fmt.Sprintf("OUT:\n%s\n", output))
		c.Done <- true
		close(c.Done)
	}

	msg := "gonna run FFmpeg"
	go runner(msg)
	return fmt.Sprintf("INFO: %s", msg)
}
