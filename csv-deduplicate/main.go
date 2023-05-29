package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	INPUT_FILE      = "./labeled_images_db.csv"
	OUTPUT_FILE     = "./labeled_images_db_cleaned.csv"
	OUTPUT_HEADLINE = "filename,label"
)

func main() {
	csvLines, err := readCSVFile(INPUT_FILE, 1)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("number of csv lines:", len(csvLines))

	deduplicatedLines := deduplicateLines(csvLines)
	log.Println("number of deduplicated lines:", len(deduplicatedLines))

	err = writeCSVFile(OUTPUT_FILE, OUTPUT_HEADLINE, deduplicatedLines)
	if err != nil {
		log.Fatalln(err)
	}
}

func readCSVFile(filename string, offset int) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	lines := [][]string{}
	for i := 0; i < offset; i++ {
		if _, err := csvReader.Read(); err != nil {
			log.Fatalln(err)
		}
	}
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lines = append(lines, rec)
	}
	return lines, nil
}

func deduplicateLines(lines [][]string) [][]string {
	dedupMap := make(map[string]string)
	for _, l := range lines {
		dedupMap[l[0]] = l[1]
	}
	dedup := make([][]string, len(dedupMap))
	i := 0
	for k, v := range dedupMap {
		dedup[i] = []string{k, v}
		i++
	}
	return dedup
}

func writeCSVFile(filename string, headline string, lines [][]string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	//write headline
	_, err = fmt.Fprintln(f, headline)
	if err != nil {
		return err
	}
	//write all data lines
	csvWriter := csv.NewWriter(f)
	csvWriter.WriteAll(lines)
	return nil
}
