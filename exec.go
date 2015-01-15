package exec

import (
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/dockpit/iowait"
)

//Cmd extends the original os/exec Command structure
type Cmd struct {
	*exec.Cmd
}

//Works like start but waits a given amount of time
//for a line from stdout that matches
//the provided regexp, if the command has an stdout or
//sterr that is not nill, it will replace the stdout with
//a io.Multiwriter that also pipes output to the stream
//that is scanned for the regex
func (c *Cmd) StartWithTimeout(to time.Duration, exp *regexp.Regexp) error {

	var err error
	var opipe io.Reader
	if c.Stdout != nil {
		var w io.Writer
		opipe, w = io.Pipe()

		c.Stdout = io.MultiWriter(w, c.Stdout)
	} else {
		opipe, err = c.StdoutPipe()
		if err != nil {
			return err
		}
	}

	var epipe io.Reader
	if c.Stderr != nil {
		var w io.Writer
		epipe, w = io.Pipe()

		c.Stderr = io.MultiWriter(w, c.Stderr)
	} else {
		epipe, err = c.StderrPipe()
		if err != nil {
			return err
		}
	}

	c.Start()

	return iowait.WaitForRegexp(io.MultiReader(epipe, opipe), exp, to)
}

//Attemps to gracefully shut down the process by first
//sending a interrupt signal and then wait the given amount
//for the process to shut down, if the process is still running
//kill it. @todo It will send a os.Interrupt signal which is not
//currently supported in windows
func (c *Cmd) StopWithTimeout(to time.Duration) error {
	exited := make(chan bool)
	go func() {
		c.Wait()
		exited <- true
	}()

	ps, err := c.expandToChildProcesses(c.Process)
	if err != nil {
		return err
	}

	//signal all (child)processes to better simulate
	//shell behaviour
	for _, p := range ps {
		err := p.Signal(os.Interrupt)
		if err != nil {
			return err
		}
	}

	select {
	case <-exited:
		return nil //process exited by itself
	case <-time.After(to):
		for _, p := range ps {
			err := p.Kill()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Command(name string, args ...string) *Cmd {
	return &Cmd{exec.Command(name, args...)}
}
