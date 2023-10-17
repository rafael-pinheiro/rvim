package cursor

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func blinkCursor(style chan lipgloss.Style) tea.Cmd {
	blink := true
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * 900)
			blink = !blink
			style <- styles[blink]
		}
	}
}

type blinkMsg struct {
	style lipgloss.Style
}

func waitForBlink(style chan lipgloss.Style) tea.Cmd {
	return func() tea.Msg {
		return blinkMsg{
			style: <-style,
		}
	}
}

type moveMsg struct {
	direction string
	amount    int
}

func Move(direction string, amount int) tea.Cmd {
	return func() tea.Msg {
		return moveMsg{
			direction: direction,
			amount:    amount,
		}
	}
}
