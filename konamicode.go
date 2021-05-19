// Copyright (C) 2021  Antonio Lassandro

// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.

// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.

// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const usage = "konamicode <file>"

type Instruction uint
type Identifier string

const (
	InstructionUp Instruction = iota
	InstructionDown
	InstructionLeft
	InstructionRight
	InstructionA
	InstructionB
	InstructionStart
	InstructionSelect
)

const (
	IdentifierUp     Identifier = "up"
	IdentifierDown              = "down"
	IdentifierLeft              = "left"
	IdentifierRight             = "right"
	IdentifierA                 = "a"
	IdentifierB                 = "b"
	IdentifierStart             = "start"
	IdentifierSelect            = "select"
)

func compile(reader io.Reader) (result []Instruction, err error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	result = make([]Instruction, 0, 256)
	branches := 0

	for scanner.Scan() {
		ident := scanner.Text()
		switch Identifier(strings.ToLower(ident)) {
		case IdentifierUp:
			result = append(result, InstructionUp)
		case IdentifierDown:
			result = append(result, InstructionDown)
		case IdentifierLeft:
			result = append(result, InstructionLeft)
		case IdentifierRight:
			result = append(result, InstructionRight)
		case IdentifierA:
			result = append(result, InstructionA)
		case IdentifierB:
			result = append(result, InstructionB)
		case IdentifierStart:
			branches++
			result = append(result, InstructionStart)
		case IdentifierSelect:
			if branches == 0 {
				result = nil
				err = errors.New("Unexpected ']'")
				return
			}

			branches--
			result = append(result, InstructionSelect)
		default:
			result = nil
			err = fmt.Errorf("Invalid identifier '%s'\n", ident)
			return
		}
	}

	if err = scanner.Err(); err != nil {
		result = nil
		return
	}

	if branches > 0 {
		result = nil
		err = fmt.Errorf("Missing closing ']' for %d branch(es)", branches)
		return
	}

	return
}

func execute(instructions []Instruction) error {
	memory := make([]byte, 30000)
	data := 0
	program := 0

	stdin := bufio.NewReader(os.Stdin)
	stack := make([]int, 0, 32)

	for program < len(instructions) {
		switch instructions[program] {
		case InstructionUp:
			memory[data]++
		case InstructionDown:
			memory[data]--
		case InstructionLeft:
			data--
		case InstructionRight:
			data++
		case InstructionA:
			if c, err := stdin.ReadByte(); err != nil && err != io.EOF {
				return err
			} else if err != io.EOF {
				memory[data] = c
			}

		case InstructionB:
			fmt.Printf("%c", memory[data])

		case InstructionStart:
			if memory[data] == 0 {
				for instructions[program] != InstructionSelect {
					program++
				}
				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, program)
			}

		case InstructionSelect:
			if memory[data] != 0 {
				program = stack[len(stack)-1]
			} else {
				stack = stack[:len(stack)-1]
			}
		}

		if program < 0 || program >= len(instructions) {
			return fmt.Errorf("Program counter (%d) out of bounds!", program)
		}

		if data < 0 || data >= len(memory) {
			return fmt.Errorf("Data pointer (%d) out of bounds!", data)
		}

		program++
	}

	return nil
}

func init() {
	exe, _ := os.Executable()
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", filepath.Base(exe)))
	log.SetOutput(os.Stderr)
}

func konamicode() int {
	var input io.Reader
	var instructions []Instruction
	var err error

	if stat, err := os.Stdin.Stat(); err != nil {
		log.Println(err)
		return 1
	} else if stat.Mode()&os.ModeCharDevice == 0 {
		input = os.Stdin
	} else {
		if len(os.Args) != 2 {
			log.Println(usage)
			return 1
		}

		var file *os.File
		file, err = os.Open(os.Args[1])

		if err != nil {
			log.Println(err)
			return 1
		}

		defer file.Close()
		input = file
	}

	if instructions, err = compile(input); err != nil {
		log.Println(err)
		return 1
	}

	if err = execute(instructions); err != nil {
		log.Println(err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(konamicode())
}
