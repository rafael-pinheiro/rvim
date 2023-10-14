package editor

import (
	"fmt"
	"rvim/pkg/buffer"
	"rvim/pkg/cursor"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	buffer buffer.Model
	cursor cursor.Model
	width  int
	height int
}

func CreateModel(filePath string) Model {
	buffer := buffer.CreateModel(filePath)

	return Model{
		buffer: buffer,
		cursor: cursor.CreateModel(&buffer),
		width:  0,
		height: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.cursor.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var separator = "▏"
var lineNumberStyle = lipgloss.NewStyle().Width(3).Align(lipgloss.Right)

func (m Model) View() string {
	output := ""

	for line, text := range m.buffer.GetText() {
		output += fmt.Sprintf(
			"%s %s %s\n",
			lineNumberStyle.Render(m.cursor.GetFormattedDistance(line)),
			separator,
			m.cursor.View(text, line),
		)

		if line == m.height-3 {
			break
		}
	}

	return output
}
