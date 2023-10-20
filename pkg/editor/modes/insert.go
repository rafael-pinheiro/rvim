package editor

import (
	"rvim/pkg/buffer"
	"rvim/pkg/cursor"

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
			return m, cursor.Move("left", 1)
		case tea.KeyRight:
			return m, cursor.Move("right", 1)
		case tea.KeyUp:
			return m, cursor.Move("up", 1)
		case tea.KeyDown:
			return m, cursor.Move("down", 1)
		case tea.KeyEscape:
			return CreateNormalMode(), nil
		case tea.KeyDelete:
			// Handle delete
		case tea.KeyRunes:
			return m, buffer.Append(0, msg.String())
		}
	}

	return m, nil
}
