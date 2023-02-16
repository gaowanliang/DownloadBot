// Package goTree create and print tree.
package goTree

import (
	"strings"
)

const (
	newLine      = "\n"
	emptySpace   = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

var maximumLength = 36

type (
	tree struct {
		text  string
		items []Tree
	}

	// Tree is tree interface
	Tree interface {
		Add(text string) Tree
		AddTree(tree Tree)
		Items() []Tree
		Text() string
		Print() string
	}

	printer struct {
	}

	// Printer is printer interface
	Printer interface {
		Print(Tree) string
	}
)

func SetMaxLength(a int) {
	maximumLength = a
}

//New returns a new GoTree.Tree
func New(text string) Tree {
	return &tree{
		text:  text,
		items: []Tree{},
	}
}

//Add adds a node to the tree
func (t *tree) Add(text string) Tree {
	n := New(text)
	t.items = append(t.items, n)
	return n
}

//AddTree adds a tree as an item
func (t *tree) AddTree(tree Tree) {
	t.items = append(t.items, tree)
}

//Text returns the node's value
func (t *tree) Text() string {
	return t.text
}

//Items returns all items in the tree
func (t *tree) Items() []Tree {
	return t.items
}

//Print returns an visual representation of the tree
func (t *tree) Print() string {
	return newPrinter().Print(t)
}

func newPrinter() Printer {
	return &printer{}
}

//Print prints a tree to a string
func (p *printer) Print(t Tree) string {
	return t.Text() + newLine + p.printItems(t.Items(), []bool{})
}

func (p *printer) printText(text string, spaces []bool, last bool) string {
	var result string
	for _, space := range spaces {
		if space {
			result += emptySpace
		} else {
			result += continueItem
		}
	}

	indicator := middleItem
	if last {
		indicator = lastItem
	}
	var newText = text
	var leftCount = maximumLength - len(indicator) // the maximum length of the line is 36 characters, but need to subtract the length of the draw table characters

	if leftCount > 0 {
		newText = ""
		runeText := []rune(text) // if you file name is Chinese,Japanese,Korean and other unicode characters, the step value of for loop  will become 2, which will affect the mod judgment, but if it is converted to the []rune, this question will be solved
		for i, val := range runeText {
			if i%leftCount == 0 && i != 0 && i != len(text)-1 { // when up to length, need to add line break
				newText += string(val) + "\n"
				continue
			}
			newText += string(val)
		}
	}
	if newText[len(newText)-1] == '\n' {
		newText = newText[:len(newText)-1]
	}
	//newText = strings.ReplaceAll(newText, "\n\n", "\n")
	var out string
	lines := strings.Split(newText, "\n")
	for i := range lines {
		text := lines[i]
		if i == 0 {
			out += result + indicator + text + newLine
			continue
		}
		if last {
			indicator = emptySpace
		} else {
			indicator = continueItem
		}
		out += result + indicator + text + newLine
		//log.Println(out)
	}

	return out
}

func (p *printer) printItems(t []Tree, spaces []bool) string {
	var result string
	for i, f := range t {
		last := i == len(t)-1
		result += p.printText(strings.ReplaceAll(f.Text(), "\n", ""), spaces, last)
		if len(f.Items()) > 0 {
			spacesChild := append(spaces, last)
			result += p.printItems(f.Items(), spacesChild)
		}
	}
	return result
}
