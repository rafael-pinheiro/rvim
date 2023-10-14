package main

import (
	"flag"
	"fmt"
	"os"
	"rvim/pkg/editor"

	tea "github.com/charmbracelet/bubbletea"
)

type rootModel struct {
	activeEditor editor.Model
}

func createModel(filePath string) rootModel {
	return rootModel{
		activeEditor: editor.CreateModel(filePath),
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.activeEditor.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.activeEditor, cmd = m.activeEditor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m rootModel) View() string {
	return m.activeEditor.View()
}

func main() {
	flag.Parse()

	filePath := flag.Arg(0)

	p := tea.NewProgram(
		createModel(filePath),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
