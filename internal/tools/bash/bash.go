package bash

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BashSession struct {
	cmd            *exec.Cmd
	stdin          io.WriteCloser
	stdout         *bufio.Reader
	stderr         *bufio.Reader
	stdinPrompt    string
	defaultTimeout time.Duration
}

type BashSessionOption func(*BashSession)

func WithTimeout(timeout time.Duration) BashSessionOption {
	return func(bs *BashSession) {
		bs.defaultTimeout = timeout
	}
}

func NewBashSession(opts ...BashSessionOption) (*BashSession, error) {
	cmd := exec.Command("bash", "--norc", "--noprofile")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	sessionId := uuid.New()
	stdInPrompt := fmt.Sprintf("'__READY_%v__'\n", sessionId)

	defaultTimeout, err := time.ParseDuration("1m")
	if err != nil {
		return nil, err
	}

	bs := &BashSession{
		cmd:            cmd,
		stdin:          stdin,
		stdout:         bufio.NewReader(stdout),
		stderr:         bufio.NewReader(stderr),
		stdinPrompt:    stdInPrompt,
		defaultTimeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(bs)
	}

	bs.stdin.Write([]byte(fmt.Sprintf("PS1=%v", stdInPrompt)))
	for {
		line, _ := bs.stdout.ReadString('\n')
		if strings.Contains(line, stdInPrompt) {
			break
		}
	}

	return bs, nil
}

func (bs *BashSession) ExecuteWithTimeout(command string, timeout time.Duration) (string, error) {
	_, err := bs.stdin.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	var output strings.Builder
	done := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		for {
			line, err := bs.stdout.ReadString('\n')
			if err != nil {
				errChan <- err
			}

			if strings.Contains(line, bs.stdinPrompt) {
				done <- output.String()
				return
			}

			output.WriteString(line)
		}
	}()

	select {
	case result := <-done:
		return strings.TrimRight(result, "\n"), nil
	case err := <-errChan:
		return "", err
	case <-time.After(timeout):
		return "", fmt.Errorf("Command timed out after %v", timeout.String())
	}
}

func (bs *BashSession) Execute(command string) (string, error) {
	return bs.ExecuteWithTimeout(command, bs.defaultTimeout)
}
