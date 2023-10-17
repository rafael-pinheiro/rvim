package cursor

import (
	"fmt"
	"rvim/pkg/buffer"
	"rvim/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	buffer         *buffer.Model
	line           int
	column         int
	direction      string
	style          lipgloss.Style
	blinkerChannel chan lipgloss.Style
}

var styles = map[bool]lipgloss.Style{
	true: lipgloss.NewStyle().
		Background(lipgloss.Color("#CCCCCC")),
	false: lipgloss.NewStyle(),
}

func CreateModel(buffer *buffer.Model) Model {
	return Model{
		buffer:         buffer,
		line:           0,
		column:         0,
		direction:      "",
		style:          styles[true],
		blinkerChannel: make(chan lipgloss.Style),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		blinkCursor(m.blinkerChannel),
		waitForBlink(m.blinkerChannel),
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

func (m *Model) goLeftBy(amount int) {
	m.column = utils.Max(0, m.column-amount)
	m.direction = "left"
	m.style = styles[true]
}

func (m *Model) goRightBy(amount int) {
	m.column = utils.Min(len(m.buffer.GetLine(m.line))-1, m.column+amount)
	m.direction = "right"
}

func (m *Model) GoUpBy(amount int) {
	m.line = utils.Max(0, m.line-amount)
	m.direction = "up"

	m.column = utils.Min(
		utils.Max(0, len(m.buffer.GetLine(m.line))-1),
		m.column,
	)

	m.style = styles[true]
}

func (m *Model) GoUp() {
	m.GoUpBy(1)
}

func (m *Model) GoDownBy(amount int) {
	m.line = utils.Min(len(m.buffer.GetText())-1, m.line+amount)
	m.direction = "down"
	m.column = utils.Min(
		utils.Max(0, len(m.buffer.GetLine(m.line))-1),
		m.column,
	)
	m.style = styles[true]
}

func (m *Model) GoDown() {
	m.GoDownBy(1)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case blinkMsg:
		m.style = msg.style
		return m, waitForBlink(m.blinkerChannel)
	case moveMsg:
		switch msg.direction {
		case "right":
			m.column = utils.Min(len(m.buffer.GetLine(m.line))-1, m.column+msg.amount)
		case "left":
			m.column = utils.Max(0, m.column-msg.amount)
		case "up":
			m.line = utils.Max(0, m.line-msg.amount)

			m.column = utils.Min(
				utils.Max(0, len(m.buffer.GetLine(m.line))-1),
				m.column,
			)
		case "down":
			m.line = utils.Min(len(m.buffer.GetText())-1, m.line+msg.amount)
			m.column = utils.Min(
				utils.Max(0, len(m.buffer.GetLine(m.line))-1),
				m.column,
			)
		default:
			panic(fmt.Sprintf("Cannot move cursor, direction \"%s\" is invalid", msg.direction))
		}

		m.direction = msg.direction
		m.style = styles[true]
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
