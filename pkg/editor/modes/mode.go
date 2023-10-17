package editor

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Mode interface {
	Update(msg tea.Msg) (Mode, tea.Cmd)
}
