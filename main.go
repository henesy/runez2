package main

import (
	"bufio"
	"container/list"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

// Encodes a rune and a position for decompression
type Pair struct {
	R rune
	P uint8
}

var chatty bool

// Naive, but revised, utf-8 text compression/decompression program
func main() {
	c := flag.Bool("c", false, "Explicit compress mode")
	d := flag.Bool("d", false, "De-compress mode")
	flag.BoolVar(&chatty, "D", false, "Chatty debug mode")
	flag.Parse()

	// Compress by default
	if (!*c) && (!*d) {
		*c = true
	}

	if *c == *d {
		fatal("err: choose one explicit mode")
	}

	// TODO - allow calling like `runez2 -c foo.txt` ;; use flag.Args

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	// Choose mode operation
	switch {
	case *c:
		Compress(in, out)
	case *d:
		Decompress(in, out)
	}
}

// Compress text to the archive format
func Compress(r *bufio.Reader, w *bufio.Writer) {
	dict := make(map[rune]uint8)
	runes := list.New()
	var first rune

	// Build table
	// i is the index position of unique runes
	for i := 0; ; {
		ru, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			fatal("err: could not read rune - ", err)
		}

		// Check for character count overflow
		if max := int(^uint8(0)); i > max {
			fatal("err: too many runes to index, stopped at: ", max)
		}

		if i <= 0 {
			// Handle the base case explicitly
			dict[ru] = uint8(i)
			i++
			first = ru
		} else {
			// 0 means the rune isn't indexed - kind of a hack
			if dict[ru] <= 0 && ru != first {
				dict[ru] = uint8(i)
				i++
			}
		}

		// Push the i-value for the rune
		runes.PushBack(dict[ru])
	}

	table := make([]rune, len(dict))

	// Iterate the dict to build the table
	for ru, i := range dict {
		if chatty {
			fmt.Fprintf(os.Stderr, "%q has index %v\n", ru, i)
		}

		table[i] = ru
	}

	// Iterate table to emit output file format
	for i := 0; i < len(table); i++ {
		if chatty {
			fmt.Fprintf(os.Stderr, "%q emitted from table\n", table[i])
		}

		// Rune
		err := binary.Write(w, binary.LittleEndian, table[i])
		if err != nil {
			fatal("err: could not write rune - ", err)
		}
	}

	// Null Rune
	err := binary.Write(w, binary.LittleEndian, rune(0))
	if err != nil {
		fatal("err: could not write rune - ", err)
	}

	// Indices
	for p := runes.Front(); p != nil; p = p.Next() {
		err := binary.Write(w, binary.LittleEndian, byte(p.Value.(uint8)))
		if err != nil {
			fatal("err: could not write position - ", err)
		}
	}

	w.Flush()
}

// Decompress the archive format to text
func Decompress(r *bufio.Reader, w *bufio.Writer) {
	var table []rune	// Store final table of runes

	table = make([]rune, 0, int(^uint8(0)))

	// Read out runes
	for {
		var ru rune

		err := binary.Read(r, binary.LittleEndian, &ru)
		if err != nil {
			fatal("err: could not read rune - ", err)
		}

		if ru == 0 {
			// \0 rune
			if chatty {
				fmt.Fprintln(os.Stderr, ">>> Hit \\0 rune")
			}
			break
		}

		table = append(table, ru)
	}

	if chatty {
		fmt.Fprintf(os.Stderr, "Table = {\n")
		for _, ru := range table {
			fmt.Fprintf(os.Stderr, "%q,\n", ru)
		}
		fmt.Fprintf(os.Stderr, "}\n")
	}

	// Read in indices and emit ouput rune
	for {
			var i uint8

			err := binary.Read(r, binary.LittleEndian, &i)
			if err != nil {
				if err == io.EOF {
					break
				}

				fatal("err: could not read index - ", err)
			}

			if int(i) > len(table) || int(i) < 0 {
				fatal("err: bad archive, out of bounds index number: ", i)
			}

			w.Write([]byte(string(table[i])))
	}

	w.Flush()
}

// Fatal - end program with an error message and newline
func fatal(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
	os.Exit(1)
}
