[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Coverage Status](https://coveralls.io/repos/github/t34-dev/go-text-replacer/badge.svg?branch=main&v=1722596450)](https://coveralls.io/github/t34-dev/go-text-replacer?branch=main&v=1722596450)
![Go Version](https://img.shields.io/badge/Go-1.22-blue?logo=go&v=1722596450)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/t34-dev/go-text-replacer?v=1722596450)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/t34-dev/go-text-replacer?sort=semver&style=flat&logo=git&logoColor=white&label=Latest%20Version&color=blue&v=1722596450)

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

Here's a simple example of how to use Go-Text-Replacer:

```go
package main

import (
	"fmt"
	"log"

	nestedreplacer "github.com/t34-dev/go-text-replacer"
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

	replacer := nestedreplacer.NewFromString(content)

	blocks := []nestedreplacer.Block{
		replacer.CreateBlockFromString("Introduction", "Getting Started with"),
		replacer.CreateBlockFromString("Programming", "Coding"),
		replacer.CreateBlockFromString("(算法思维)", "(Algorithmic Thinking)"),
	}

	result, err := replacer.Enter(blocks)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(result))
}
```

## API Reference

### Types

- `Block`: Represents a replacement block with Start, End, and Txt fields.
- `Position`: Represents the position of found text with Start and End fields.

### Functions

- `New(content []byte) *nestedReplacer`: Creates a new instance of nestedReplacer with the given content.
- `NewFromString(content string) *nestedReplacer`: Creates a new instance of nestedReplacer from a string.

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
