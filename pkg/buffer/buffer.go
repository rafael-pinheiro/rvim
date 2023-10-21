package buffer

import (
	"bufio"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type lineDescriptor struct {
	pieces []PieceDescriptor
}

type PieceDescriptor struct {
	source string
	offset int64
	length int64
}

type Model struct {
	original         *os.File
	append           *os.File
	appendOffset     int64
	pieceDescriptors []PieceDescriptor
	lines            int
}

func CreateModel(filePath string) Model {
	originalFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	appendFile, err := os.CreateTemp("", "append-")

	if err != nil {
		fmt.Println("could not create temporary file:", err)
		os.Exit(1)
	}

	stat, err := originalFile.Stat()

	if err != nil {
		fmt.Println("could not read file:", err)
		os.Exit(1)
	}

	return Model{
		original: originalFile,
		append:   appendFile,
		pieceDescriptors: []PieceDescriptor{{
			source: "original",
			offset: 0,
			length: stat.Size(),
		}},
		lines: -1,
	}
}

func (m Model) Update(msg tea.Msg, position int) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case appendMsg:
		length := int64(len(msg.text))
		m.append.WriteString(msg.text)
		m.pieceDescriptors = append([]PieceDescriptor{{
			source: "append",
			offset: m.appendOffset,
			length: length,
		}}, m.pieceDescriptors...)
		m.appendOffset += length
	}

	return m, nil
}

func (m Model) GetLines(offset int, length int) []string {
	scanner := m.newScanner()
	var currentLine int

	for currentLine < offset {
		scanner.Scan()
		currentLine++
	}

	output := make([]string, 0, length)
	currentLine = 0

	for currentLine < length {
		scanner.Scan()
		output = append(output, scanner.Text())
		currentLine++
	}

	return output
}

func (m Model) newScanner() *bufio.Scanner {
	var readers []io.Reader

	for _, piece := range m.pieceDescriptors {
		reader := m.original
		if piece.source == "append" {
			reader = m.append
		}

		readers = append(readers, io.NewSectionReader(reader, piece.offset, piece.length))
	}
	reader := io.MultiReader(readers...)

	return bufio.NewScanner(reader)
}

func (m Model) countLines() int {
	scanner := m.newScanner()
	var lines int

	for scanner.Scan() {
		lines++
	}

	return lines
}

func (m *Model) CountLines() int {
	if m.lines == -1 {
		m.lines = m.countLines()
	}

	return m.lines
}
