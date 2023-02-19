package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"crg.eti.br/go/config"
)

type config struct {
	InputFile  string `json:"input_file" cfg:"i" cfgDefault:"-"`
	OutputFile string `json:"output_file" cfg:"o" cfgDefault:"-"`
	Comma      string `json:"comma" cfg:"comma" cfgDefault:","`
}

func closer(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config{}

	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	csvFile := os.Stdin

	if cfg.InputFile != "-" {
		csvFile, err = os.Open(cfg.InputFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer closer(csvFile)
	}

	r := csv.NewReader(csvFile)
	r.Comma = ';'
	// r.FieldsPerRecord

	outFile := os.Stdout
	if cfg.OutputFile != "-" {
		outFile, err = os.Create(cfg.OutputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	writer := bufio.NewWriter(outFile)

	h := []string{}
	i := 0

	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println(err)

			return
		}

		if i == 0 {
			h = append(h, record[0])
			i++

			continue
		}

		_, err = writer.WriteString(fmt.Sprintf("%s = %s\n", h[0], record[0]))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
