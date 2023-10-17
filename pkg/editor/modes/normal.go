package editor

import (
	"rvim/pkg/cursor"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type NormalMode struct {
	name     string
	repeater string
}

func CreateNormalMode() NormalMode {
	return NormalMode{
		name:     "NORMAL",
		repeater: "",
	}
}

func (m NormalMode) Update(msg tea.Msg) (Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			return m, cursor.Move("left", m.getRepeater())
		case "l", "right":
			return m, cursor.Move("right", m.getRepeater())
		case "k", "up":
			return m, cursor.Move("up", m.getRepeater())
		case "j", "down":
			return m, cursor.Move("down", m.getRepeater())
		case "i":
			return CreateInsertMode(), nil
		case "0":
			if m.repeater != "" {
				m.repeater += "0"
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.repeater += msg.String()
		}
	}

	return m, nil
}

func (m *NormalMode) getRepeater() int {
	repeater, err := strconv.Atoi(m.repeater)

	if err != nil {
		return 1
	}

	m.repeater = ""

	return repeater
}
