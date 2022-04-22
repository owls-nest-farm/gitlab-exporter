package main

import (
	"encoding/csv"
	"errors"
	"io/fs"
	"os"
)

func main() {
	f, err := os.Open("export.csv")
	if err != nil {
		panic(err)
	}
	csv := csv.NewReader(f)
	records, err := csv.ReadAll()
	if err != nil {
		panic(err)
	}
	exporter := NewExporter(records, "https://gitlab.com/api/v4/")
	err = exporter.Export()
	if err != nil {
		panic(err)
	}
	exporter.Compress()
	if _, err := os.Stat(exporter.TmpDir); !errors.Is(err, fs.ErrNotExist) {
		if err := os.RemoveAll(exporter.TmpDir); err != nil {
			panic(err)
		}
	}
}
