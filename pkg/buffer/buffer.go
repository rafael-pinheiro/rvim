package buffer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"rvim/pkg/commands"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	original         *os.File
	append           *os.File
	appendOffset     int64
	pieceDescriptors []PieceDescriptor
	currentPiece     int
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

func (m Model) offsetFromPosition(line int, column int) int64 {
	_, offset := m.newScannerAtLine(line)

	return int64(offset + column)
}

func (m Model) appendAt(text string, offset int64) (Model, tea.Cmd) {
	pieces := m.pieceDescriptors
	m.append.WriteString(text)

	var pieceIndex int
	var piece PieceDescriptor
	var currentOffset int64
	var innerOffset int64
	length := int64(len(text))

	for i, p := range pieces {
		if offset > currentOffset && offset <= currentOffset+p.length {
			pieceIndex = i
			piece = p
			innerOffset = offset - currentOffset
			break
		}
		currentOffset += p.length
	}

	// At the end of the piece
	if innerOffset == piece.length {
		// Continuous typing grows the current piece
		if pieceIndex == m.currentPiece {
			m.pieceDescriptors[pieceIndex].length += length
			m.appendOffset += length
			return m, commands.Move("right", int(length))
		}

		// Creates a new piece if it's not the current piece
		newPiece := PieceDescriptor{
			source: "append",
			offset: m.appendOffset,
			length: length,
		}
		m.currentPiece = pieceIndex + 1
		m.pieceDescriptors = append(
			m.pieceDescriptors[:pieceIndex+1],
			append(
				[]PieceDescriptor{newPiece},
				m.pieceDescriptors[pieceIndex+1:]...,
			)...,
		)

		m.appendOffset += length

		return m, commands.Move("right", int(length))
	}

	// Splits a piece in two and inserts new piece between
	newPiece := PieceDescriptor{
		source: "append",
		offset: m.appendOffset,
		length: length,
	}
	m.currentPiece = pieceIndex + 1

	m.pieceDescriptors = append(
		m.pieceDescriptors[:pieceIndex],
		append(
			[]PieceDescriptor{
				// First split of the current piece
				{
					source: piece.source,
					offset: piece.offset,
					length: innerOffset,
				},
				newPiece,
				// Second split of the current piece
				{
					source: piece.source,
					offset: piece.offset + innerOffset,
					length: piece.length - innerOffset,
				},
			},
			m.pieceDescriptors[pieceIndex:]...,
		)...,
	)

	m.appendOffset += length

	return m, commands.Move("right", 1)
}

func (m Model) Update(msg tea.Msg, line int, column int) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commands.AppendMsg:
		return m.appendAt(msg.Text, m.offsetFromPosition(line, column))
	}

	return m, nil
}

func (m Model) GetLines(offset int, length int) []string {
	scanner, _ := m.newScannerAtLine(offset)

	output := make([]string, 0, length)
	var currentLine int

	for currentLine < length {
		scanner.Scan()
		output = append(output, scanner.Text())
		currentLine++
	}

	return output
}

func (m Model) newScanner() *bufio.Scanner {
	var readers []io.Reader
	pieces := m.pieceDescriptors

	for _, piece := range pieces {
		reader := m.original
		if piece.source == "append" {
			reader = m.append
		}

		readers = append(readers, io.NewSectionReader(reader, piece.offset, piece.length))
	}
	reader := io.MultiReader(readers...)

	return bufio.NewScanner(reader)
}

func (m Model) newScannerAtLine(targetLine int) (*bufio.Scanner, int) {
	var currentLine int
	var discardedBytes int
	scanner := m.newScanner()

	for currentLine < targetLine {
		scanner.Scan()
		discardedBytes += len(scanner.Bytes())
		currentLine++
	}

	return scanner, discardedBytes
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
