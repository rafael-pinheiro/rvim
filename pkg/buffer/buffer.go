package buffer

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type lineDescriptor struct {
	offset int
	length int
}

type Model struct {
	original        io.ReaderAt
	size            int
	lineDescriptors []lineDescriptor
}

func CreateModel(filePath string) Model {
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	stat, err := file.Stat()

	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	lineDescriptors := []lineDescriptor{}

	scanner := bufio.NewScanner(file)
	offset := 0
	for scanner.Scan() {
		lineLength := len(scanner.Bytes())
		lineDescriptors = append(lineDescriptors, lineDescriptor{
			offset: offset,
			length: lineLength,
		})
		offset += lineLength + 1
	}

	return Model{
		original:        file,
		size:            int(stat.Size()),
		lineDescriptors: lineDescriptors,
	}
}

func (m Model) Lines() int {
	return len(m.lineDescriptors)
}

func (m Model) GetLine(line int) string {
	lineDescriptor := m.lineDescriptors[line]
	reader := io.NewSectionReader(m.original, int64(lineDescriptor.offset), int64(lineDescriptor.length))
	result := make([]byte, lineDescriptor.length)
	_, err := reader.Read(result)

	if err != nil {
		panic(err)
	}

	return string(result)
}
