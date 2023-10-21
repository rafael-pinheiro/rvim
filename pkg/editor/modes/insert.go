package editor

import (
	"rvim/pkg/commands"

	tea "github.com/charmbracelet/bubbletea"
)

type InsertMode struct {
	name string
}

func CreateInsertMode() InsertMode {
	return InsertMode{
		name: "INSERT",
	}
}

func (m InsertMode) Update(msg tea.Msg) (Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			return m, commands.Move("left", 1)
		case tea.KeyRight:
			return m, commands.Move("right", 1)
		case tea.KeyUp:
			return m, commands.Move("up", 1)
		case tea.KeyDown:
			return m, commands.Move("down", 1)
		case tea.KeyEscape:
			return CreateNormalMode(), nil
		case tea.KeyDelete:
			// Handle delete
		case tea.KeyRunes:
			return m, commands.Append(msg.String())
		}
	}

	return m, nil
}
