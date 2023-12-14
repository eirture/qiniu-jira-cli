package cmdutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type askOptions struct {
	DefaultValue string
}

func (opt *askOptions) Apply(options ...AskOption) *askOptions {
	for _, o := range options {
		o(opt)
	}
	return opt
}

type AskOption func(opt *askOptions)

func WithAskDefaultValue(defaultValue string) AskOption {
	return func(o *askOptions) {
		o.DefaultValue = defaultValue
	}
}

func Ask(prompt string, options ...AskOption) (answer string, err error) {
	opt := (&askOptions{}).Apply(options...)
	prompt = strings.TrimSpace(prompt)
	if opt.DefaultValue != "" {
		prompt = prompt + fmt.Sprintf(" [%s]", opt.DefaultValue)
	}
	_, err = fmt.Fprint(os.Stdout, prompt+": ")
	if err != nil {
		return
	}
	reader := bufio.NewReader(os.Stdin)
	answer, err = reader.ReadString('\n')
	if err != nil {
		return
	}
	answer = strings.Replace(answer, "\n", "", -1)
	if answer == "" {
		answer = opt.DefaultValue
	}
	return
}

func AskVar(prompt string, v *string, options ...AskOption) (err error) {
	if *v != "" {
		options = append([]AskOption{WithAskDefaultValue(*v)}, options...)
	}
	answer, err := Ask(prompt, options...)
	if err != nil {
		return
	}
	*v = answer
	return
}
