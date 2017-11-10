package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/manifoldco/promptui"

	"github.com/enmand/quarid-go/pkg/vm/python"
)

const EOF = "EOF"

func main() {
	logrus.Info("Loading Python script")
	py := python.New()

	file := ""
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	if file == "" {
		for {
			code, err := repl(">")
			if err.Error() == EOF {
				return
			}
			if err != nil && err != promptui.ErrEOF {
				panic(fmt.Sprintf("python error: %s", err))
			}

			err = py.Run(code)
			if err != nil {
				fmt.Printf("error running python code: %s", err)
			}
		}
	}
}

func repl(prompt string) (string, error) {
	code, err := replPrompt(prompt)
	if err == promptui.ErrEOF {
		return code, err // return code and EOF
	}
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(code, ":") {
		for {
			newCode, err := repl("#")
			if err == promptui.ErrEOF {
				return code, nil
			}
			if err != nil {
				return "", err
			}

			if code == "" {
				return code, nil
			}

			code = fmt.Sprintf("%s\n%s", code, newCode)
		}
	}

	return code, nil
}

func replPrompt(prompt string) (string, error) {
	p := promptui.Prompt{
		Label:     fmt.Sprintf("quarid %s", prompt),
		IsVimMode: true,
		Validate: func(code string) error {
			vm := python.New()
			return vm.Run(code)
		},
	}

	return p.Run()
}
