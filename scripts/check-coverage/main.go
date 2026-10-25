package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	profile := flag.String("profile", "coverage.out", "Go coverage profile")
	minimum := flag.Float64("min", 80, "minimum statement coverage percentage")
	flag.Parse()
	if err := checkFile(*profile, *minimum); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func checkFile(path string, minimum float64) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	coverage, err := readCoverage(bytes.NewReader(data))
	if err != nil {
		return err
	}
	fmt.Printf("Statement coverage: %.2f%%; required: %.2f%%\n", coverage, minimum)
	if coverage < minimum {
		return fmt.Errorf("coverage %.2f%% is below %.2f%%", coverage, minimum)
	}
	return nil
}

func readCoverage(input io.Reader) (float64, error) {
	scanner := bufio.NewScanner(input)
	blocks := make(map[string]coverageBlock)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		name, block, err := parseBlock(line)
		if err != nil {
			return 0, err
		}
		block.covered = block.covered || blocks[name].covered
		blocks[name] = block
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	var covered, total int
	for _, block := range blocks {
		total += block.statements
		if block.covered {
			covered += block.statements
		}
	}
	if total == 0 {
		return 0, errors.New("coverage profile contains no statements")
	}
	return 100 * float64(covered) / float64(total), nil
}

type coverageBlock struct {
	statements int
	covered    bool
}

func parseBlock(line string) (string, coverageBlock, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return "", coverageBlock{}, fmt.Errorf("invalid coverage block: %q", line)
	}
	statements, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", coverageBlock{}, err
	}
	count, err := strconv.Atoi(fields[2])
	return fields[0], coverageBlock{statements: statements, covered: count > 0}, err
}
