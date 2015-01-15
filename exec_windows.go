package exec

import (
	"os"
)

func (c *Cmd) expandToChildProcesses(p *os.Process) ([]*os.Process, error) {
	return []*os.Process{p}, nil
}
