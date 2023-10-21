package cursor

import (
	"fmt"
	"rvim/pkg/buffer"
	"rvim/pkg/commands"
	"rvim/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	line           int
	column         int
	direction      string
	style          lipgloss.Style
	blinkerChannel chan lipgloss.Style
}

func CreateModel() Model {
	return Model{
		line:           0,
		column:         0,
		direction:      "",
		blinkerChannel: make(chan lipgloss.Style),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		commands.BlinkCursor(m.blinkerChannel),
		commands.WaitForBlink(m.blinkerChannel),
	)
}

func (m Model) GetDistance(line int) int {
	return line - m.line
}

func (m Model) GetPosition() (int, int) {
	return m.line, m.column
}

func (m Model) GetDirection() string {
	return m.direction
}

func (m Model) getLineLength(buffer buffer.Model) int {
	return len(buffer.GetLines(m.line, 1)[0])
}

func (m Model) Update(msg tea.Msg, buffer buffer.Model) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commands.BlinkMsg:
		m.style = msg.Style
		return m, commands.WaitForBlink(m.blinkerChannel)
	case commands.MoveMsg:
		switch msg.Direction {
		case "right":
			m.column = utils.Min(m.getLineLength(buffer)-1, m.column+msg.Amount)
		case "left":
			m.column = utils.Max(0, m.column-msg.Amount)
		case "up":
			m.line = utils.Max(0, m.line-msg.Amount)

			m.column = utils.Min(
				utils.Max(0, m.getLineLength(buffer)-1),
				m.column,
			)
		case "down":
			m.line = utils.Min(buffer.CountLines()-1, m.line+msg.Amount)
			m.column = utils.Min(
				utils.Max(0, m.getLineLength(buffer)-1),
				m.column,
			)
		default:
			panic(fmt.Sprintf("Cannot move cursor, direction \"%s\" is invalid", msg.Direction))
		}

		m.direction = msg.Direction
		// m.style = styles[true]
	}

	return m, nil
}

func (m Model) View(text string, line int) string {
	if line != m.line {
		return text
	}
	if len(text) == 0 {
		return m.style.Render(" ")
	}

	return fmt.Sprintf(
		"%s%s%s",
		text[:m.column],
		m.style.Render(string(text[m.column])),
		text[m.column+1:],
	)
}
