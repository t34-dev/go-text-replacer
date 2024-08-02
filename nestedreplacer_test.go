package nestedreplacer

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEnter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		blocks      []Block
		expected    string
		expectError bool
	}{
		{
			name:     "Basic single block replacement",
			content:  "Hello, world!",
			blocks:   []Block{{Start: 7, End: 12, Txt: []byte("Go")}},
			expected: "Hello, Go!",
		},
		{
			name:    "Multiple block replacements",
			content: "The quick brown fox jumps over the lazy dog",
			blocks: []Block{
				{Start: 4, End: 9, Txt: []byte("slow")},
				{Start: 16, End: 19, Txt: []byte("cat")},
				{Start: 35, End: 39, Txt: []byte("active")},
			},
			expected: "The slow brown cat jumps over the active dog",
		},
		{
			name:    "Overlapping blocks",
			content: "abcdefghijkl",
			blocks: []Block{
				{Start: 3, End: 6, Txt: []byte("XXX")},
				{Start: 5, End: 8, Txt: []byte("YYY")},
			},
			expected:    "",
			expectError: true,
		},
		{
			name:     "Replacement at the beginning",
			content:  "Start here",
			blocks:   []Block{{Start: 0, End: 5, Txt: []byte("Begin")}},
			expected: "Begin here",
		},
		{
			name:     "Replacement at the end",
			content:  "End there",
			blocks:   []Block{{Start: 4, End: 9, Txt: []byte("here")}},
			expected: "End here",
		},
		{
			name:     "Replace entire content",
			content:  "Replace me",
			blocks:   []Block{{Start: 0, End: 10, Txt: []byte("New content")}},
			expected: "New content",
		},
		{
			name:     "Empty content",
			content:  "",
			blocks:   []Block{{Start: 0, End: 0, Txt: []byte("Something")}},
			expected: "",
		},
		{
			name:     "Empty blocks",
			content:  "Don't change me",
			blocks:   []Block{},
			expected: "Don't change me",
		},
		{
			name:    "Blocks outside content (should be ignored)",
			content: "Short",
			blocks: []Block{
				{Start: 10, End: 15, Txt: []byte("Ignored")},
				{Start: 5, End: 5, Txt: []byte("Finish")},
				{Start: 0, End: 5, Txt: []byte("Long")},
			},
			expected: "LongFinish",
		},
		{
			name:    "Blocks with negative positions (should be ignored)",
			content: "Negative",
			blocks: []Block{
				{Start: -5, End: 0, Txt: []byte("Ignored")},
				{Start: 0, End: 4, Txt: []byte("Posi")},
			},
			expected: "Positive",
		},
		{
			name:    "Blocks where End <= Start (should be ignored)",
			content: "Incorrect",
			blocks: []Block{
				{Start: 5, End: 3, Txt: []byte("Ignored")},
				{Start: 0, End: 2, Txt: []byte("Co")},
			},
			expected:    "",
			expectError: true,
		},
		{
			name:    "Unsorted blocks (should be sorted internally)",
			content: "First Second Third",
			blocks: []Block{
				{Start: 13, End: 18, Txt: []byte("Last")},
				{Start: 0, End: 5, Txt: []byte("Start")},
			},
			expected: "Start Second Last",
		},
		{
			name:    "Multibyte characters",
			content: "Hello, 世界!",
			blocks: []Block{
				{Start: 7, End: -1, Txt: []byte("мир!")},
			},
			expected: "Hello, мир!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nr := New([]byte(tt.content))
			result, err := nr.Enter(tt.blocks)

			if tt.expectError && err == nil {
				t.Errorf("Expected an error, but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && !bytes.Equal(result, []byte(tt.expected)) {
				t.Errorf("Enter() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

func TestFindAllPositions(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		searchText string
		expected   []Position
	}{
		{
			name:       "Text at the beginning",
			content:    "Hello, world!",
			searchText: "Hello",
			expected:   []Position{{Start: 0, End: 5}},
		},
		{
			name:       "Text at the end",
			content:    "Hello, world!",
			searchText: "world!",
			expected:   []Position{{Start: 7, End: 13}},
		},
		{
			name:       "Text in the middle",
			content:    "The quick brown fox",
			searchText: "quick",
			expected:   []Position{{Start: 4, End: 9}},
		},
		{
			name:       "Text not found",
			content:    "Hello, world!",
			searchText: "goodbye",
			expected:   nil,
		},
		{
			name:       "Multiple occurrences",
			content:    "The quick brown fox is quick",
			searchText: "quick",
			expected:   []Position{{Start: 4, End: 9}, {Start: 23, End: 28}},
		},
		{
			name:       "Empty content",
			content:    "",
			searchText: "something",
			expected:   nil,
		},
		{
			name:       "Empty search text",
			content:    "Not empty",
			searchText: "",
			expected:   nil,
		},
		{
			name:       "Multibyte characters",
			content:    "Hello, 世界! 世界 is world.",
			searchText: "世界",
			expected:   []Position{{Start: 7, End: 13}, {Start: 15, End: 21}},
		},
		{
			name:       "Search text longer than content",
			content:    "Short",
			searchText: "This is too long",
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nr := New([]byte(tt.content))
			positions := nr.FindAllPositions([]byte(tt.searchText))

			if !reflect.DeepEqual(positions, tt.expected) {
				t.Errorf("FindAllPositions() = %v, want %v", positions, tt.expected)
			}
		})
	}
}

func TestFindFirstPosition(t *testing.T) {
	nr := New([]byte("Hello, world! Hello again!"))

	tests := []struct {
		name       string
		searchText string
		startIndex int
		expected   *Position
	}{
		{
			name:       "Find first occurrence",
			searchText: "Hello",
			startIndex: 0,
			expected:   &Position{Start: 0, End: 5},
		},
		{
			name:       "Find second occurrence",
			searchText: "Hello",
			startIndex: 1,
			expected:   &Position{Start: 14, End: 19},
		},
		{
			name:       "Text not found",
			searchText: "Goodbye",
			startIndex: 0,
			expected:   nil,
		},
		{
			name:       "Start index out of bounds",
			searchText: "Hello",
			startIndex: 100,
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			position := nr.FindFirstPosition([]byte(tt.searchText), tt.startIndex)

			if !reflect.DeepEqual(position, tt.expected) {
				t.Errorf("FindFirstPosition() = %v, want %v", position, tt.expected)
			}
		})
	}
}

func TestFindLastPosition(t *testing.T) {
	nr := New([]byte("Hello, world! Hello again!"))

	tests := []struct {
		name       string
		searchText string
		startIndex int
		expected   *Position
	}{
		{
			name:       "Find last occurrence",
			searchText: "Hello",
			startIndex: -1,
			expected:   &Position{Start: 14, End: 19},
		},
		{
			name:       "Find first occurrence",
			searchText: "Hello",
			startIndex: 13,
			expected:   &Position{Start: 0, End: 5},
		},
		{
			name:       "Text not found",
			searchText: "Goodbye",
			startIndex: -1,
			expected:   nil,
		},
		{
			name:       "Start index out of bounds",
			searchText: "Hello",
			startIndex: 100,
			expected:   &Position{Start: 14, End: 19},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			position := nr.FindLastPosition([]byte(tt.searchText), tt.startIndex)

			if !reflect.DeepEqual(position, tt.expected) {
				t.Errorf("FindLastPosition() = %v, want %v", position, tt.expected)
			}
		})
	}
}

func TestCreateBlock(t *testing.T) {
	nr := New([]byte("The quick brown fox"))

	tests := []struct {
		name     string
		find     string
		replace  string
		expected Block
	}{
		{
			name:     "Create block for existing text",
			find:     "quick",
			replace:  "slow",
			expected: Block{Start: 4, End: 9, Txt: []byte("slow")},
		},
		{
			name:     "Create block for non-existing text",
			find:     "lazy",
			replace:  "energetic",
			expected: Block{Start: -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := nr.CreateBlock([]byte(tt.find), []byte(tt.replace))

			if !reflect.DeepEqual(block, tt.expected) {
				t.Errorf("CreateBlock() = %v, want %v", block, tt.expected)
			}
		})
	}
}

func TestRuneToBytePosition(t *testing.T) {
	nr := New([]byte("Hello, 世界!"))

	tests := []struct {
		name          string
		runeStart     int
		runeEnd       int
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "ASCII characters",
			runeStart:     0,
			runeEnd:       5,
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "Include multibyte characters",
			runeStart:     7,
			runeEnd:       9,
			expectedStart: 7,
			expectedEnd:   13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteStart, byteEnd := nr.RuneToBytePosition(tt.runeStart, tt.runeEnd)

			if byteStart != tt.expectedStart || byteEnd != tt.expectedEnd {
				t.Errorf("RuneToBytePosition() = (%v, %v), want (%v, %v)", byteStart, byteEnd, tt.expectedStart, tt.expectedEnd)
			}
		})
	}
}

func TestByteToRunePosition(t *testing.T) {
	nr := New([]byte("Hello, 世界!"))

	tests := []struct {
		name          string
		byteStart     int
		byteEnd       int
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "ASCII characters",
			byteStart:     0,
			byteEnd:       5,
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "Include multibyte characters",
			byteStart:     7,
			byteEnd:       11,
			expectedStart: 7,
			expectedEnd:   9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runeStart, runeEnd := nr.ByteToRunePosition(tt.byteStart, tt.byteEnd)

			if runeStart != tt.expectedStart || runeEnd != tt.expectedEnd {
				t.Errorf("ByteToRunePosition() = (%v, %v), want (%v, %v)", runeStart, runeEnd, tt.expectedStart, tt.expectedEnd)
			}
		})
	}
}

func TestCreateBlockFromString(t *testing.T) {
	nr := NewFromString("Hello, 世界!")

	tests := []struct {
		name     string
		find     string
		replace  string
		expected Block
	}{
		{
			name:     "Create block for ASCII text",
			find:     "Hello",
			replace:  "Bonjour",
			expected: Block{Start: 0, End: 5, Txt: []byte("Bonjour")},
		},
		{
			name:     "Create block for multibyte text",
			find:     "世界",
			replace:  "мир",
			expected: Block{Start: 7, End: 13, Txt: []byte("мир")},
		},
		{
			name:     "Create block for non-existing text",
			find:     "Goodbye",
			replace:  "Adieu",
			expected: Block{Start: -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := nr.CreateBlockFromString(tt.find, tt.replace)

			if block.Start != tt.expected.Start {
				t.Errorf("CreateBlockFromString() Start = %v, want %v", block.Start, tt.expected.Start)
			}
			if block.End != tt.expected.End {
				t.Errorf("CreateBlockFromString() End = %v, want %v", block.End, tt.expected.End)
			}
			if !bytes.Equal(block.Txt, tt.expected.Txt) {
				t.Errorf("CreateBlockFromString() Txt = %v, want %v", block.Txt, tt.expected.Txt)
			}
		})
	}
}
