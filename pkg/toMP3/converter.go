package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	tplDIR_FAKE = "files"
	tplDIR_REAL = "output"
	tplEXT_NEW  = ".mp3"
)

type conv struct {
	Done chan bool
	Err  error
	Info []string
	pth  string
}

func New(pth string) *conv {
	c := conv{
		Done: make(chan bool),
		pth:  pth,
	}
	_ = dir.getApp()
	return &c
}

func (c *conv) Check() bool {
	var (
		fin      = c.pth
		setError = func(num int) bool {
			c.Err = fmt.Errorf("file path error, for '%s' [%d]", c.pth, num)
			return false
		}
	)

	if fin == "" {
		return setError(1)
	}

	lst := strings.Split(fin, "/")
	if len(lst) < 2 {
		return setError(2)
	}
	val := strings.TrimSpace(lst[0])
	if val == "." || val == "" { // remove first element if it '.' or empty
		lst = lst[1:]
	}
	lst[0] = tplDIR_REAL // replace user entry for static files with real dir name

	if err := dir.getApp(); err != nil {
		c.Info = append(c.Info, err.Error())
	} else {
		lst = append([]string{dir.app}, lst...) // prepend exec. path

	}
	fin = filepath.Join(lst...) // make file path as string

	if _, err := os.Stat(fin); errors.Is(err, os.ErrNotExist) {
		return setError(3)
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
		out := filepath.Join(pth, (fl[0:len(fl)-len(filepath.Ext(fl))] + tplEXT_NEW))

		ffmpeg, err := dir.getFFmpeg()  // need to check error
		if ffmpeg == "" && err != nil { // because could be no FFmpeg installed
			c.Info = append(c.Info, err.Error())
			c.Err = err
			c.done()
			return
		}

		cmd := exec.Command(ffmpeg, []string{
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
		c.done()
	}

	msg := "gonna run FFmpeg"
	go runner(msg)
	return fmt.Sprintf("INFO: %s", msg)
}

func (c *conv) done() {
	c.Done <- true
	close(c.Done)
}

var (
	dir defPath
)

type defPath struct {
	app    string
	ffmpeg string
}

func (d *defPath) getApp() error {
	if d.app != "" {
		return nil
	}

	fl, err := os.Executable()
	if err != nil {
		return err
	}

	d.app = filepath.Dir(fl)
	return nil
}

func (d *defPath) getFFmpeg() (string, error) {
	if d.ffmpeg != "" {
		return d.ffmpeg, nil
	}

	// FIXME: this is not for POSIX
	cmd := exec.Command("which", "ffmpeg")
	output, _ := cmd.CombinedOutput()

	err := cmd.Run()

	// FIXME: do not work with '\r\n' (Win)
	d.ffmpeg = strings.TrimSuffix(string(output), "\n")

	if d.ffmpeg == "" {
		return "", fmt.Errorf("cannot find FFmpeg")
	}

	return d.ffmpeg, err
}
