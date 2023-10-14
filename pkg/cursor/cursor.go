package cursor

import (
	"fmt"
	"rvim/pkg/buffer"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	buffer         *buffer.Model
	line           int
	column         int
	repeater       string
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
		repeater:       "",
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

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func (m *Model) getRepeater() int {
	repeater, err := strconv.Atoi(m.repeater)

	if err != nil {
		return 1
	}

	m.repeater = ""

	return repeater
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.column = max(0, m.column-m.getRepeater())
			m.direction = "left"
			m.style = styles[true]
		case "l", "right":
			m.column = min(len(m.buffer.GetLine(m.line))-1, m.column+m.getRepeater())
			m.direction = "right"
			m.style = styles[true]
		case "k", "up":
			m.line = max(0, m.line-m.getRepeater())
			m.direction = "up"
			m.column = min(
				max(0, len(m.buffer.GetLine(m.line))-1),
				m.column,
			)
			m.style = styles[true]
		case "j", "down":
			m.line = min(len(m.buffer.GetText())-1, m.line+m.getRepeater())
			m.direction = "down"
			m.column = min(
				max(0, len(m.buffer.GetLine(m.line))-1),
				m.column,
			)
			m.style = styles[true]
		case "0":
			if m.repeater != "" {
				m.repeater += "0"
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.repeater += msg.String()
		}

	case blinkMsg:
		m.style = msg.style
		return m, waitForBlink(m.blinkerChannel)
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

func blinkCursor(style chan lipgloss.Style) tea.Cmd {
	blink := true
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * 900)
			blink = !blink
			style <- styles[blink]
		}
	}
}

type blinkMsg struct {
	style lipgloss.Style
}

func waitForBlink(style chan lipgloss.Style) tea.Cmd {
	return func() tea.Msg {
		return blinkMsg{
			style: <-style,
		}
	}
}
