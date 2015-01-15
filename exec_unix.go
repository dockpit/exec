// +build darwin linux

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// Add low level options that start the new process into a new group
func (c *Cmd) startInNewGroup() error {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	return nil
}

// Retrieve child processes of the given process using pgrep and return its children and the
// process itself
func (c *Cmd) expandToChildProcesses(p *os.Process) ([]*os.Process, error) {

	//get the group id of the new process
	gid, err := syscall.Getpgid(p.Pid)
	if err != nil {
		return nil, fmt.Errorf("(unix) Failed to get pgid: %s", err)
	}

	//get all processes in the group using system command (works on darwin since 10.8)
	pgrep := exec.Command("pgrep", "-g", strconv.Itoa(gid))
	output, err := pgrep.Output()
	if err != nil {
		return nil, fmt.Errorf("(unix) Failed to get pgrep output: %s", err)
	}

	//get children pids
	ps := []*os.Process{p}
	cpids := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cpids := range cpids {
		cpid, err := strconv.Atoi(cpids)
		if err != nil {
			return nil, fmt.Errorf("(unix) Failed to convert pid from '%s' to int: %s", cpids, err)
		}

		cp, err := os.FindProcess(cpid)
		if err != nil {
			return nil, fmt.Errorf("(unix) Failed to find pid '%d': %s", cpid, err)
		}

		ps = append(ps, cp)
	}

	return ps, nil
}
