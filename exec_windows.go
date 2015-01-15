// +build windows

package exec

import (
	"os"
	"syscall"
)

func (c *Cmd) startInNewGroup() error {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
	}

	return nil
}

func (c *Cmd) expandToChildProcesses(p *os.Process) ([]*os.Process, error) {
	return []*os.Process{p}, nil
}
