package prompter

type Prompt struct {
	// template Templates
	ow Overwriter
}

func New() *Prompt {
	return &Prompt{}
}

func (p *Prompt) Bool(question string, prefill bool) bool {
	return Bool(question, prefill)
}

func (p *Prompt) Password(label string, allowBlank bool) (string, error) {
	return Password(label, allowBlank)
}

func (p *Prompt) Select(question string, hideHelp bool, options ...string) (string, error) {
	return Select(question, hideHelp, options...)
}

func (p *Prompt) Number(question string, prefill int, allowBlank bool) (int, error) {
	return Number(question, prefill, allowBlank)
}

func (p *Prompt) Text(question string, prefill string, allowBlank bool) (string, error) {
	return Text(question, prefill, allowBlank)
}

func (p *Prompt) NewOverwriter(question string, hideHelp bool) {
	p.ow, _ = NewOverwriter(question, hideHelp)
}

func (p *Prompt) OverwriteBool(question string, hideHelp, original bool, prompt func() bool) (bool, error) {
	return p.ow.Bool(question, hideHelp, original, prompt)
}

func (p *Prompt) OverwriteSelect(question string, hideHelp bool, original string, prompt func() string) (string, error) {
	return p.ow.Select(question, hideHelp, original, prompt)
}

func (p *Prompt) OverwriteText(question string, hideHelp bool, original string, prompt func() string) (string, error) {
	return p.ow.Text(question, hideHelp, original, prompt)
}

func (p *Prompt) OverwritePassword(question string, hideHelp bool, original string, prompt func() string) (string, error) {
	return p.ow.Password(question, hideHelp, original, prompt)
}
