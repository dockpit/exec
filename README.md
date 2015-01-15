# exec
[![GoDoc](http://godoc.org/github.com/dockpit/exec?status.png)](http://godoc.org/github.com/dockpit/exec)

Enhances os/exec Cmd with timed-out start and stop, usefull for long running processes.

TODO
====
Investigate windows support with examples @https://github.com/tgulacsi/go/tree/master/proc:

```
func procAttrSetGroup(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func isGroupLeader(c *exec.Cmd) bool {
	return c.SysProcAttr != nil &&
		c.SysProcAttr.CreationFlags&syscall.CREATE_NEW_PROCESS_GROUP > 0
}

// Pkill kills the process with the given pid
func Pkill(pid int) error {
	return exec.Command("taskkill", "/f", "/pid", strconv.Itoa(pid)).Run()
}

// GroupKill kills the process group lead by the given pid
func GroupKill(pid int) error {
	return exec.Command("taskkill", "/f", "/t", strconv.Itoa(pid)).Run()
}
```
