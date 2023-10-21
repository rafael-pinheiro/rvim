package commands

import tea "github.com/charmbracelet/bubbletea"

type AppendMsg struct {
	Text string
}

func Append(text string) tea.Cmd {
	return func() tea.Msg {
		return AppendMsg{
			Text: text,
		}
	}
}
