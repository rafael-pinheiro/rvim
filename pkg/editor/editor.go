package editor

import (
	"fmt"
	"rvim/pkg/buffer"
	"rvim/pkg/cursor"
	editor "rvim/pkg/editor/modes"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	buffer  buffer.Model
	cursor  cursor.Model
	mode    editor.Mode
	topLine int
	width   int
	height  int
}

func CreateModel(filePath string) Model {
	buffer := buffer.CreateModel(filePath)

	return Model{
		buffer:  buffer,
		cursor:  cursor.CreateModel(),
		mode:    editor.CreateNormalMode(),
		topLine: 0,
		width:   10,
		height:  10,
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

	m.cursor, cmd = m.cursor.Update(msg, m.buffer)
	cmds = append(cmds, cmd)

	m.mode, cmd = m.mode.Update(msg)
	cmds = append(cmds, cmd)

	_, column := m.cursor.GetPosition()
	m.buffer, cmd = m.buffer.Update(msg, column)
	cmds = append(cmds, cmd)

	m.adjustViewPort()

	return m, tea.Batch(cmds...)
}

var separator = "▏"
var lineNumberStyle = lipgloss.NewStyle().Width(3).Align(lipgloss.Right)

func (m Model) viewLineNumber(line int) string {
	lineDistance := m.cursor.GetDistance(line)
	lineNumber := fmt.Sprint(line + 1)
	if lineDistance < 0 {
		lineNumber = fmt.Sprint(lineDistance)
	} else if lineDistance > 0 {
		lineNumber = fmt.Sprintf("+%d", lineDistance)
	}

	return lineNumberStyle.Render(lineNumber)
}

func (m *Model) adjustViewPort() {
	line, _ := m.cursor.GetPosition()
	direction := m.cursor.GetDirection()

	switch direction {
	case "up":
		if line < m.topLine {
			m.topLine = line
		}
	case "down":
		if line > m.topLine+m.height-2 {
			m.topLine = line - m.height + 2
		}
	}
}

func (m Model) View() string {
	output := ""
	lines := m.buffer.GetLines(m.topLine, m.height-1)

	for index, line := range lines {
		output += fmt.Sprintf(
			"%s %s %s\n",
			m.viewLineNumber(m.topLine+index),
			separator,
			m.cursor.View(line, m.topLine+index),
		)
	}

	return output
}
