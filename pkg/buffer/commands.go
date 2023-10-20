package buffer

import tea "github.com/charmbracelet/bubbletea"

type appendMsg struct {
	offset int64
	text   string
}

func Append(offset int64, text string) tea.Cmd {
	return func() tea.Msg {
		return appendMsg{
			offset: offset,
			text:   text,
		}
	}
}
