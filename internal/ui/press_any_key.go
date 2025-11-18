package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// PressAnyKeyDialog is a Bubble Tea model for waiting for any key press.
type PressAnyKeyDialog struct {
	Message     string
	Pressed     bool
	StyleConfig StyleConfig
}

var _ tea.Model = (*PressAnyKeyDialog)(nil)

// NewPressAnyKeyDialog creates a new press key dialog.
func NewPressAnyKeyDialog(message string) PressAnyKeyDialog {
	return PressAnyKeyDialog{
		Message:     message,
		StyleConfig: DefaultStyleConfig(),
	}
}

func (d PressAnyKeyDialog) Init() tea.Cmd {
	return nil
}

func (d PressAnyKeyDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() != "" {
			d.Pressed = true
			return d, tea.Quit
		}
	}

	return d, nil
}

func (d PressAnyKeyDialog) View() string {
	if d.Pressed {
		return ""
	}

	return d.StyleConfig.HelpStyle.Render(d.Message) + "\n"
}

// Run is a convenience method to run the press key dialog.
func (d PressAnyKeyDialog) Run() error {
	p := tea.NewProgram(d)
	_, err := p.Run()
	return err
}

// WaitForKey displays a message and waits for any key press.
func WaitForKey(message string) error {
	dialog := NewPressAnyKeyDialog(message)
	return dialog.Run()
}
