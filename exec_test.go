package exec_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/dockpit/exec"
	"github.com/stretchr/testify/assert"
)

func compileDaemon(t *testing.T) {
	//compile a example daemon that we can stop
	bcmd := exec.Command("go", "build", "-o=/tmp/daemon", filepath.Join(".example", "main.go"))
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	err := bcmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartTimesOut(t *testing.T) {
	cmd := exec.Command("go", "version")

	err := cmd.StartWithTimeout(time.Millisecond*10, regexp.MustCompile(`bogus`))
	assert.Error(t, err)
}

func TestStart(t *testing.T) {
	cmd := exec.Command("go", "env")

	err := cmd.StartWithTimeout(time.Second, regexp.MustCompile(`GOPATH`))
	assert.NoError(t, err)
}

func TestStartStop(t *testing.T) {
	compileDaemon(t)

	cmd := exec.Command("/tmp/daemon")
	err := cmd.StartWithTimeout(time.Second, regexp.MustCompile(`Starting Goji`))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = cmd.StopWithTimeout(time.Second)
		assert.NoError(t, err)
	}()
}

func TestStartStopWithDefinedStd(t *testing.T) {
	compileDaemon(t)

	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("/tmp/daemon")
	cmd.Stderr = buf
	cmd.Stdout = buf

	err := cmd.StartWithTimeout(time.Second, regexp.MustCompile(`Starting Goji`))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = cmd.StopWithTimeout(time.Second)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "gracefully")
	}()
}
