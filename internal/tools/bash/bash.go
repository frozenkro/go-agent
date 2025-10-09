// Package bash provides a persistent bash session to be invoked as a tool by llm agents
// Current implementation is specific to anthropic spec:
// https://anthropic.mintlify.app/en/docs/agents-and-tools/tool-use/bash-tool

package bash

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	toolschema "github.com/frozenkro/go-agent/models/anthropic/tool_schema"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

const (
	BUFFER_SIZE      int           = 1024
	BUFFER_POLL_RATE time.Duration = time.Millisecond * 10
)

type BashTool struct {
	bs *BashSession
}

type BashSession struct {
	tty            tty
	prompt         string
	defaultTimeout time.Duration
}

type tty interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	SetReadDeadline(time.Time) error
	Close() error
}

type BashSessionOption func(*BashSession)

func WithTimeout(timeout time.Duration) BashSessionOption {
	return func(bs *BashSession) {
		bs.defaultTimeout = timeout
	}
}

func NewBashSession(opts ...BashSessionOption) (*BashSession, error) {
	cmd := exec.Command("bash", "--norc", "--noprofile", "-i")

	f, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	sessionId := uuid.New()
	prompt := fmt.Sprintf("__READY_%v__", sessionId)

	defaultTimeout, err := time.ParseDuration("1m")
	if err != nil {
		return nil, err
	}

	bs := &BashSession{
		tty:            f,
		prompt:         prompt,
		defaultTimeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(bs)
	}

	command := fmt.Sprintf("PS1=%v", prompt)
	bs.Execute(command)

	return bs, nil
}

func (bs *BashSession) ExecuteWithTimeout(command string, timeout time.Duration) (string, error) {
	bs.sendCommand(command)

	done := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		output, err := bs.getResponse(command)
		if err != nil {
			errChan <- err
		} else {
			done <- output
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
	var p toolschema.BashToolInput
	err := mapstructure.Decode(params, &p)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for BashTool: '%v'", params)
	}

	if t.bs != nil && p.Restart {
		t.bs.Deinit()
	}

	if t.bs == nil || p.Restart {
		var err error
		t.bs, err = NewBashSession()

		if err != nil {
			return "", fmt.Errorf("Error initializing Bash Session: %w", err)
		}
	}

	if p.Command != "" {
		return t.bs.Execute(p.Command)
	}

	return "", nil
}

func (bs *BashSession) getResponse(command string) (string, error) {

	buffer := make([]byte, BUFFER_SIZE)
	accumulated := ""

	for {
		iterationTime := time.Now()
		bs.tty.SetReadDeadline(iterationTime.Add(BUFFER_POLL_RATE))
		n, err := bs.tty.Read(buffer)
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			return "", err
		}

		if n > 0 {
			chunk := string(buffer[:n])
			accumulated += chunk
		}

		if bs.isAtPrompt(accumulated) {
			result := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(accumulated), bs.prompt), command)
			return result, nil
		} else {
			time.Sleep(BUFFER_POLL_RATE)
		}

	}
}

func (bs *BashSession) isAtPrompt(accumulated string) bool {
	cleanStr := strings.TrimSpace(accumulated)
	endsInPrompt := strings.HasSuffix(cleanStr, bs.prompt)
	endsInPromptInit := strings.HasSuffix(cleanStr, fmt.Sprintf("PS1=%v", bs.prompt))
	return endsInPrompt && !endsInPromptInit
}

func (bs *BashSession) sendCommand(c string) {
	bs.tty.Write([]byte(c))
	bs.tty.Write([]byte("\n"))
}

func (bs *BashSession) Deinit() {
	if bs.tty != nil {
		// First write an EOT to indicate end of `bash` command
		bs.tty.Write([]byte{4})
		time.Sleep(time.Millisecond * 100)

		// Close file descriptor for character device
		bs.tty.Close()
	}
}

// TODO Repurpose original implementation into tests
// func main() {
// 	c := exec.Command("bash", "-i")
// 	f, err := pty.Start(c)
// 	if err != nil {
// 		panic(err)
// 	}

// 	command := fmt.Sprintf("PS1='%v'", PROMPT)
// 	sendCommand(f, command)
// 	fmt.Print(getResponse(f, command))

// 	command = "ls -al"
// 	sendCommand(f, command)
// 	fmt.Print(getResponse(f, command))

// 	command = "pwd"
// 	sendCommand(f, command)
// 	fmt.Print(getResponse(f, command))

// 	f.Write([]byte{4}) // EOT

// 	fmt.Println("\nCompleted..")

// }
