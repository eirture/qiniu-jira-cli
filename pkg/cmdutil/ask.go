package cmdutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type askOptions struct {
	DefaultValue string
	isPassword   bool
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

func WithAskPassword(is bool) AskOption {
	return func(o *askOptions) {
		o.isPassword = is
	}
}

func Ask(prompt string, options ...AskOption) (answer string, err error) {
	opt := (&askOptions{}).Apply(options...)
	prompt = strings.TrimSpace(prompt)
	if opt.DefaultValue != "" {
		prompt = prompt + fmt.Sprintf(" [%s]", asteriskStrIf(opt.DefaultValue, opt.isPassword))
	}
	if opt.isPassword {
		_, err = fmt.Fprint(os.Stdout, prompt+": ")
		if err != nil {
			return
		}
		var pwd []byte
		pwd, err = term.ReadPassword(syscall.Stdin)
		if err != nil {
			return
		}
		answer = string(pwd)
		fmt.Println()
	} else {
		_, err = fmt.Fprint(os.Stdout, prompt+": ")
		if err != nil {
			return
		}
		reader := bufio.NewReader(os.Stdin)
		answer, err = reader.ReadString('\n')
		if err != nil {
			return
		}
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

func asteriskStrIf(str string, condition bool) string {
	if condition {
		return AsteriskStr(str)
	}
	return str
}

func AsteriskStr(str string) string {
	length := len(str)
	l, r := 0, length
	if length >= 5 {
		l += 2
		r -= 2
	} else if length >= 3 {
		l += 1
		r -= 1
	} else if length > 1 {
		l += 1
	}
	return str[:l] + strings.Repeat("*", r-l) + str[r:]
}
