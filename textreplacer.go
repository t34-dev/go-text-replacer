// Copyright 2024 The Go-Text-Replacer Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package textreplacer

import (
	"bytes"
	"fmt"
	"sort"
	"unicode/utf8"
)

// Block represents a structure for storing information about a replacement block.
type Block struct {
	Start int
	End   int
	Txt   []byte
}

// Position represents the position of found text.
type Position struct {
	Start int
	End   int
}

// textreplacer is a structure containing the original content.
type textreplacer struct {
	content []byte
}

// New creates a new instance of textreplacer with the given content.
func New(content []byte) *textreplacer {
	return &textreplacer{
		content: content,
	}
}

// NewFromString creates a new instance of textreplacer from a string.
func NewFromString(content string) *textreplacer {
	return &textreplacer{
		content: []byte(content),
	}
}

// Enter applies replacement blocks to the original content.
func (n *textreplacer) Enter(blocks []Block) ([]byte, error) {
	content := n.content
	if len(content) == 0 {
		return nil, nil
	}
	if len(blocks) == 0 {
		return content, nil
	}

	// Create a slice of pointers to blocks
	sortedBlocks := make([]*Block, 0, len(blocks))
	for i := range blocks {
		block := &blocks[i]
		if block.Start < 0 || block.Start > len(content) {
			continue
		}
		if block.End < 0 || block.End > len(content) {
			block.End = len(content)
		}
		if block.End < block.Start {
			return nil, fmt.Errorf("range error: start [%d] >= end [%d]", block.Start, block.End)
		}
		sortedBlocks = append(sortedBlocks, block)
	}
	if len(sortedBlocks) == 0 {
		return content, nil
	}

	// Sort the slice of pointers
	sort.Slice(sortedBlocks, func(i, j int) bool {
		return sortedBlocks[i].Start < sortedBlocks[j].Start
	})

	// Check for overlap
	for i := 1; i < len(sortedBlocks); i++ {
		if sortedBlocks[i].Start < sortedBlocks[i-1].End {
			return nil, fmt.Errorf("overlap error: block %d [%d:%d] overlaps with block %d [%d:%d]",
				i-1, sortedBlocks[i-1].Start, sortedBlocks[i-1].End,
				i, sortedBlocks[i].Start, sortedBlocks[i].End)
		}
	}

	var result []byte
	lastIndex := 0

	for _, block := range sortedBlocks {
		// Add text before the current block
		if block.Start > lastIndex {
			result = append(result, content[lastIndex:block.Start]...)
		}

		// Add replacement text
		result = append(result, block.Txt...)

		lastIndex = block.End
	}

	// Add remaining text
	if lastIndex < len(content) {
		result = append(result, content[lastIndex:]...)
	}

	return result, nil
}

// FindAllPositions finds all positions of the given text in the content.
func (n *textreplacer) FindAllPositions(text []byte) []Position {
	if len(text) == 0 {
		return nil
	}
	var positions []Position

	for i := 0; i <= len(n.content)-len(text); i++ {
		if bytes.Equal(n.content[i:i+len(text)], text) {
			positions = append(positions, Position{Start: i, End: i + len(text)})
		}
	}

	if len(positions) == 0 {
		return nil
	}

	return positions
}

// FindFirstPosition finds the first position of the given text in the content, starting from the specified index.
func (n *textreplacer) FindFirstPosition(text []byte, startIndex int) *Position {
	if len(text) == 0 {
		return nil
	}
	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex > len(n.content) {
		return nil
	}

	idx := bytes.Index(n.content[startIndex:], text)
	if idx == -1 {
		return nil
	}
	start := startIndex + idx
	return &Position{Start: start, End: start + len(text)}
}

// FindLastPosition finds the last position of the given text in the content, starting from the specified index from the end.
func (n *textreplacer) FindLastPosition(text []byte, startIndex int) *Position {
	if len(text) == 0 {
		return nil
	}
	if startIndex < 0 || startIndex >= len(n.content) {
		startIndex = len(n.content) - 1
	}

	idx := bytes.LastIndex(n.content[:startIndex+1], text)
	if idx == -1 {
		return nil
	}
	return &Position{Start: idx, End: idx + len(text)}
}

// CreateBlock creates a block using the given find and replacement text.
func (n *textreplacer) CreateBlock(find, txt []byte) Block {
	position := n.FindFirstPosition(find, 0)
	if position == nil {
		return Block{
			Start: -1,
		}
	}
	return Block{
		Start: position.Start,
		End:   position.End,
		Txt:   txt,
	}
}

// RuneToBytePosition converts a position in runes to a position in bytes.
func (n *textreplacer) RuneToBytePosition(runeStart, runeEnd int) (byteStart, byteEnd int) {
	byteStart = 0
	for i := 0; i < runeStart && byteStart < len(n.content); i++ {
		_, size := utf8.DecodeRune(n.content[byteStart:])
		byteStart += size
	}

	byteEnd = byteStart
	for i := runeStart; i < runeEnd && byteEnd < len(n.content); i++ {
		_, size := utf8.DecodeRune(n.content[byteEnd:])
		byteEnd += size
	}

	return
}

// ByteToRunePosition converts a position in bytes to a position in runes.
func (n *textreplacer) ByteToRunePosition(byteStart, byteEnd int) (runeStart, runeEnd int) {
	runeStart = 0
	for i := 0; i < byteStart; {
		_, size := utf8.DecodeRune(n.content[i:])
		i += size
		runeStart++
	}

	runeEnd = runeStart
	for i := byteStart; i < byteEnd; {
		_, size := utf8.DecodeRune(n.content[i:])
		i += size
		runeEnd++
	}

	return
}

// CreateBlockFromString creates a block using rune positions.
func (n *textreplacer) CreateBlockFromString(findRunes, txtRunes string) Block {
	find := []byte(findRunes)
	txt := []byte(txtRunes)
	position := n.FindFirstPosition(find, 0)
	if position == nil {
		return Block{
			Start: -1,
		}
	}
	runeStart, runeEnd := n.ByteToRunePosition(position.Start, position.End)
	byteStart, byteEnd := n.RuneToBytePosition(runeStart, runeEnd)
	return Block{
		Start: byteStart,
		End:   byteEnd,
		Txt:   txt,
	}
}
