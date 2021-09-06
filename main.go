package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/adlawren/go-mal/ast"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO: add ability to catch ^C, to terminate long-running command
	// ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		reader, err := ast.NewReader(scanner.Text())
		if err != nil {
			log.Error().Err(err).Msg("Failed to initialize tokenizer")
		}

		rootNode, err := reader.ParseAst()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse AST")
		}

		fmt.Println(rootNode.String())
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Scanner error")
	}
}
