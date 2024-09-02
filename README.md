[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Coverage Status](https://coveralls.io/repos/github/t34-dev/go-text-replacer/badge.svg?branch=main&ver=1723898449)](https://coveralls.io/github/t34-dev/go-text-replacer?branch=main&ver=1723898449)
![Go Version](https://img.shields.io/badge/Go-1.22-blue?logo=go&ver=1723898449)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/t34-dev/go-text-replacer?ver=1723898449)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/t34-dev/go-text-replacer?sort=semver&style=flat&logo=git&logoColor=white&label=Latest%20Version&color=blue&ver=1723898449)

# Go-Text-Replacer

Go-Text-Replacer is a Go package that provides functionality for performing nested replacements in text content. It allows you to replace multiple sections of text while maintaining the original structure and handling overlapping replacements.

## Features

- Replace multiple sections of text in a single operation
- Handle overlapping replacements
- Support for both byte slices and strings
- Find all occurrences of a given text
- Find first or last occurrence of a given text
- Convert between rune and byte positions
- Create replacement blocks from strings or byte slices

## Installation

To install Go-Text-Replacer, use `go get`:

```
go get github.com/t34-dev/go-text-replacer
```

## Usage

Here's a simple example of how to use Go-Text-Replacer with strings:

```go
package main

import (
	"fmt"
	textreplacer "github.com/t34-dev/go-text-replacer"
	"log"
)

func main() {
	content := `
Chapter 1: Introduction to Programming (编程简介)

    1.1 What is Programming (什么是编程)
        Programming is the art of creating instructions for computers.
        It includes many aspects such as:
            - Algorithmic thinking (算法思维)
            - Block structures
            - Programming language syntax
`

	replacer := textreplacer.NewFromString(content)

	blocks := []textreplacer.Block{
		replacer.CreateBlockFromString("Introduction", "Getting Started with"),
		replacer.CreateBlockFromString("Programming", "Coding"),
		replacer.CreateBlockFromString("(算法思维)", "(Algorithmic Thinking)"),
	}

	result, err := replacer.Enter(blocks)
	if err != nil {
		log.Fatalln(err)
	}

	// output:
	//Chapter 1: Getting Started with to Coding (编程简介)
	//
	//1.1 What is Programming (什么是编程)
	//Programming is the art of creating instructions for computers.
	//	It includes many aspects such as:
	//- Algorithmic thinking (Algorithmic Thinking)
	//- Block structures
	//- Programming language syntax
	
	fmt.Println(string(result))
}

```

And here's an example using byte slices:

```go
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
	// analyze the file
	replacer := textreplacer.New([]byte(content))

	// if we want to find the byte range for the searched word
	position := replacer.FindFirstPosition([]byte("Introduction"), 0)
	// output: 'Introduction' found at positions 12 to 24
	if position != nil {
		fmt.Printf("'Introduction' found at positions %d to %d\n", position.Start, position.End)
	} else {
		log.Fatalln("'Introduction' not found")
	}

	blocks := []textreplacer.Block{
		// if we know exactly the byte numbers to replace
		// output: 'Introduction' found at positions 12 to 24
		{
			Start: 12,
			End:   24,
			Txt:   []byte("What_stuck a bit"),
		},
		// if we don't know the byte numbers
		replacer.CreateBlock([]byte("to"), []byte("is your")),
		replacer.CreateBlockFromString("Programming", "[用户名]"),
		replacer.CreateBlockFromString("(算法思维)", "\n            - 另一个要点"),
	}
	// perform replacement
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

	// output:
	// 1.1 What is Programming (什么是编程)
	// Programming is the art of creating instructions for computers.
	// It includes many aspects such as:
	// - Algorithmic thinking
	// - 另一个要点
	// - END structures
	// - Programming language syntax
	fmt.Println(string(result))
}
```

## API Reference

### Types

- `Block`: Represents a replacement block with Start, End, and Txt fields.
- `Position`: Represents the position of found text with Start and End fields.

### Functions

- `New(content []byte) *textreplacer`: Creates a new instance of textreplacer with the given content.
- `NewFromString(content string) *textreplacer`: Creates a new instance of textreplacer from a string.

### Methods

- `Enter(blocks []Block) ([]byte, error)`: Applies replacement blocks to the original content.
- `FindAllPositions(text []byte) []Position`: Finds all positions of the given text in the content.
- `FindFirstPosition(text []byte, startIndex int) *Position`: Finds the first position of the given text in the content, starting from the specified index.
- `FindLastPosition(text []byte, startIndex int) *Position`: Finds the last position of the given text in the content, starting from the specified index from the end.
- `CreateBlock(find, txt []byte) Block`: Creates a block using the given find and replacement text.
- `RuneToBytePosition(runeStart, runeEnd int) (byteStart, byteEnd int)`: Converts a position in runes to a position in bytes.
- `ByteToRunePosition(byteStart, byteEnd int) (runeStart, runeEnd int)`: Converts a position in bytes to a position in runes.
- `CreateBlockFromString(findRunes, txtRunes string) Block`: Creates a block using rune positions.

## License

This project is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

Developed with ❤️ by [T34](https://github.com/t34-dev)
