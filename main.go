package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/gosidekick/goconfig"
)

type config struct {
	InputFile  string `json:"input_file" cfg:"i" cfgDefault:"-"`
	OutputFile string `json:"output_file" cfg:"o" cfgDefault:"-"`
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
		defer csvFile.Close()
	}

	r := csv.NewReader(csvFile)
	r.Comma = ';'

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
		fmt.Printf("%s = %s\n", h[0], record[0])
	}
}
