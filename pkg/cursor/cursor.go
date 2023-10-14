package cursor

import (
	"fmt"
	"rvim/pkg/buffer"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	buffer         *buffer.Model
	line           int
	column         int
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

func (m Model) GetFormattedDistance(line int) string {
	lineDistance := m.GetDistance(line)
	lineNumber := fmt.Sprint(line)
	if lineDistance < 0 {
		lineNumber = fmt.Sprint(lineDistance)
	} else if lineDistance > 0 {
		lineNumber = fmt.Sprintf("+%d", lineDistance)
	}

	return lineNumber
}

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			if m.column > 0 {
				m.column--
			}
			m.style = styles[true]
		case "l", "right":
			if m.column < len(m.buffer.GetLine(m.line))-1 {
				m.column++
			}
			m.style = styles[true]
		case "k", "up":
			if m.line > 0 {
				m.line--

				lineLength := max(0, len(m.buffer.GetLine(m.line))-1)

				if lineLength == 0 {
					m.column = 0
				}

				if m.column > lineLength {
					m.column = lineLength
				}
			}
			m.style = styles[true]
		case "j", "down":

			if m.line < len(m.buffer.GetText())-1 {
				m.line++
				lineLength := max(0, len(m.buffer.GetLine(m.line))-1)

				if lineLength == 0 {
					m.column = 0
				}

				if m.column > lineLength {
					m.column = lineLength
				}
			}
			m.style = styles[true]
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
