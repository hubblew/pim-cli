package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SpinnerDialog struct {
	Spinner spinner.Model
	Text    string
	done    bool
	err     error
}

var _ tea.Model = (*SpinnerDialog)(nil)

type doneMsg struct{}
type errorMsg error

func NewSpinnerDialog(text string) SpinnerDialog {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = DefaultStyleConfig().PromptStyle

	return SpinnerDialog{
		Spinner: s,
		Text:    text,
	}
}

func (s SpinnerDialog) Init() tea.Cmd {
	return s.Spinner.Tick
}

func (s SpinnerDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		}
	case doneMsg:
		s.done = true
		return s, tea.Quit
	case errorMsg:
		s.done = true
		s.err = msg
		return s, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		s.Spinner, cmd = s.Spinner.Update(msg)
		return s, cmd
	}
	return s, nil
}

func (s SpinnerDialog) View() string {
	if s.err != nil {
		return "Error: " + s.err.Error() + "\n"
	}
	if s.done {
		return ""
	}

	return s.Spinner.View() + " " + s.Text + "\n\n"
}

// RunWithSpinner displays a Spinner while executing the provided function.
func RunWithSpinner(text string, fn func() error) error {
	dialog := NewSpinnerDialog(text)

	p := tea.NewProgram(dialog)

	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		errChan <- fn()
		p.Send(doneMsg{})
	}()

	finalModel, err := p.Run()
	if err != nil {
		cancel()
		return err
	}

	// Check if the spinner was cancelled (e.g., Ctrl+C)
	if result, ok := finalModel.(SpinnerDialog); ok && !result.done {
		cancel()
		// Wait briefly for the goroutine to finish or timeout
		select {
		case fnErr := <-errChan:
			return fnErr
		case <-time.After(100 * time.Millisecond):
			// Goroutine is still running, but we're exiting
			return context.Canceled
		}
	}

	// Normal completion - wait for the function result
	select {
	case fnErr := <-errChan:
		return fnErr
	case <-ctx.Done():
		return ctx.Err()
	}
}
