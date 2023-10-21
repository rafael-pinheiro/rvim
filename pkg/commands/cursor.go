package commands

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var styles = map[bool]lipgloss.Style{
	true: lipgloss.NewStyle().
		Background(lipgloss.Color("#CCCCCC")),
	false: lipgloss.NewStyle(),
}

func BlinkCursor(style chan lipgloss.Style) tea.Cmd {
	blink := true
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * 900)
			blink = !blink
			style <- styles[blink]
		}
	}
}

type BlinkMsg struct {
	Style lipgloss.Style
}

func WaitForBlink(style chan lipgloss.Style) tea.Cmd {
	return func() tea.Msg {
		return BlinkMsg{
			Style: <-style,
		}
	}
}

type MoveMsg struct {
	Direction string
	Amount    int
}

func Move(direction string, amount int) tea.Cmd {
	return func() tea.Msg {
		return MoveMsg{
			Direction: direction,
			Amount:    amount,
		}
	}
}
