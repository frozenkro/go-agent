package bash

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	toolModels "github.com/frozenkro/go-agent/models/anthropic/tools"
	"github.com/google/uuid"
)

type BashTool struct {
	bs *BashSession
}

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
	stdInPrompt := fmt.Sprintf("'__READY_%v__'", sessionId)

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

	_, err = bs.Execute(fmt.Sprintf("PS1=%v\n", stdInPrompt))
	if err != nil {
		return nil, err
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

func (t BashTool) Invoke(params any) (string, error) {
	bytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Unable to marshal invoke params for BashTool: '%v'", params)
	}

	var input toolModels.BashToolInput
	err = json.Unmarshal(bytes, &input)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal invoke params for BashTool: '%v'", params)
	}

	if t.bs == nil || input.Restart {
		var err error
		t.bs, err = NewBashSession()

		if err != nil {
			return "", fmt.Errorf("Error initializing Bash Session: %w", err)
		}
	}

	if input.Command != "" {
		return t.bs.Execute(input.Command)
	}

	return "", nil
}
