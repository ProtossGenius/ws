package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/ProtossGenius/pglang/snreader"
	"golang.org/x/crypto/ssh/terminal"
)

func parseGrammer(c <-chan snreader.ProductItf) {
	for {
		it := (<-c).(*lex_pgl.LexProduct)
		typ, value := it.ProductType(), it.Value
		if value == "\t" {
			fmt.Println("haha")
		}
		if typ == int(lex_pgl.PGLA_PRODUCT_HAN) {
			fmt.Println("Warn")
		}

	}
}

func read() {
	var char rune
	var err error
	sm := lex_pgl.NewLexAnalysiser()
	go parseGrammer(sm.GetResultChan())
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(fd, state)
	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err = reader.ReadRune()
		if char != '\t' {
			fmt.Printf("%c", char)
		}
		if err = sm.Read(&lex_pgl.PglaInput{Char: char}); err != nil {
			log.Fatal("when parse lex error : ", err)
		}
	}
}

func main() {

}
