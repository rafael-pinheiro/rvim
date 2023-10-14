package buffer

import (
	"bufio"
	"fmt"
	"os"
)

type Model struct {
	text []string
}

func CreateModel(filePath string) Model {
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return Model{
		text: fileLines,
	}
}

func (m Model) GetText() []string {
	return m.text
}

func (m Model) GetLine(line int) string {
	return m.text[line]
}
