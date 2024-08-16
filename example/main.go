package main

import (
	"fmt"
	textreplacer "github.com/t34-dev/go-text-replacer"
	"log"
)

var content = `
Chapter 1: Introduction to Programming (编程简介)

    1.1 What is Programming (什么是编程)
        Programming is the art of creating instructions for computers.
        It includes many aspects such as:
            - Algorithmic thinking (算法思维)
            - Block structures
            - Programming language syntax
`

func main() {
	replacer := textreplacer.NewFromString(content)
	position := replacer.FindFirstPosition([]byte("Introduction"), 0)
	if position != nil {
		fmt.Printf("'Introduction' found at positions %d to %d\n", position.Start, position.End)
	} else {
		log.Fatalln("'Introduction' not found")
	}

	blocks := []textreplacer.Block{
		{
			Start: 12,
			End:   24,
			Txt:   []byte("What_stuck a bit"),
		},
		replacer.CreateBlock([]byte("to"), []byte("is your")),
		replacer.CreateBlockFromString("Programming", "[用户名]"),
		replacer.CreateBlockFromString("(算法思维)", "\n            - 另一个要点"),
	}
	result, err := replacer.Enter(blocks)

	if err != nil {
		log.Fatalln(err)
	}

	replacer = textreplacer.New(result)
	result, err = replacer.Enter([]textreplacer.Block{
		replacer.CreateBlock([]byte("Block"), []byte("END")),
	})

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(result))
}
