package buffer

import (
	"bufio"
	"fmt"
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
	original        *os.File
	append          *os.File
	appendOffset    int64
	lineDescriptors []lineDescriptor
}

func CreateModel(filePath string) Model {
	var originalFile, appendFile *os.File
	var err error

	originalFile, err = os.Open(filePath)

	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	appendFile, err = os.CreateTemp("", "append-")

	if err != nil {
		fmt.Println("could not create temporary file:", err)
		os.Exit(1)
	}

	lineDescriptors := []lineDescriptor{}

	scanner := bufio.NewScanner(originalFile)
	var offset int64 = 0
	for scanner.Scan() {
		lineLength := int64(len(scanner.Bytes()))
		lineDescriptors = append(lineDescriptors, lineDescriptor{
			pieces: []PieceDescriptor{{
				source: "original",
				offset: offset,
				length: lineLength,
			}},
		})
		offset += lineLength + 1
	}

	return Model{
		original:        originalFile,
		append:          appendFile,
		lineDescriptors: lineDescriptors,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case appendMsg:
		length := int64(len(msg.text))
		m.append.Seek(m.appendOffset, 0)
		m.append.WriteString(msg.text)
		m.lineDescriptors[0].pieces = append([]PieceDescriptor{{
			source: "append",
			offset: m.appendOffset,
			length: length,
		}}, m.lineDescriptors[0].pieces...)
		m.appendOffset += length
	}

	return m, nil
}

func (m Model) Lines() int {
	return len(m.lineDescriptors)
}

func (m Model) GetLine(line int) string {
	lineDescriptor := m.lineDescriptors[line]
	var output []byte

	for _, piece := range lineDescriptor.pieces {
		reader := m.original
		if piece.source == "append" {
			reader = m.append
		}
		reader.Seek(piece.offset, 0)

		a := make([]byte, piece.length)
		if _, err := reader.Read(a); err != nil {
			fmt.Println(m.append.Name())
			fmt.Println(piece)
			os.Exit(1)
		}
		output = append(output, a...)
	}

	return string(output)
}
