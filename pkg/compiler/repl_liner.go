package compiler

import (
	"github.com/glycerine/liner"
)

type Prompter struct {
	prompt   string
	prompter *liner.State
	origMode liner.ModeApplier
	rawMode  liner.ModeApplier
}

func NewPrompter(prompt string) *Prompter {
	origMode, err := liner.TerminalMode()
	if err != nil {
		panic(err)
	}

	p := &Prompter{
		prompt:   prompt,
		prompter: liner.NewLiner(),
		origMode: origMode,
	}

	rawMode, err := liner.TerminalMode()
	if err != nil {
		panic(err)
	}
	p.rawMode = rawMode

	p.prompter.SetCtrlCAborts(false)

	return p
}

func (p *Prompter) Close() {
	defer p.prompter.Close()
}

func (p *Prompter) Getline(prompt *string) (line string, err error) {
	applyErr := p.rawMode.ApplyMode()
	if applyErr != nil {
		panic(applyErr)
	}
	defer func() {
		applyErr := p.origMode.ApplyMode()
		if applyErr != nil {
			panic(applyErr)
		}
	}()

	if prompt == nil {
		line, err = p.prompter.Prompt(p.prompt)
	} else {
		line, err = p.prompter.Prompt(*prompt)
	}
	if err == nil {
		p.prompter.AppendHistory(line)
		return line, nil
	}
	return "", err
}
