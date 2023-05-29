package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

const (
	CSV_FILE                 = "input/housing.csv"
	CSV_FILE_SKIP_FIRST_ROWS = 1
	NETWORK_OUTPUT_FILE      = "output/housing_network"
)

func main() {
	csvData, err := preprocess()
	if err != nil {
		panic(fmt.Sprint("preprocessing failed: ", err.Error()))
	}
	fmt.Println("number of usable rows in csv file: ", len(csvData))
	csvData.Shuffle()

	trainingData, testingData := csvData.Split(0.99)
	fmt.Println("number of usable rows in training set: ", len(trainingData))
	fmt.Println("number of usable rows in testing set: ", len(testingData))

	n := deep.NewNeural(&deep.Config{
		Inputs:     7,
		Layout:     []int{8, 4, 1},
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeBinary,
		Weight:     deep.NewNormal(1.0, 0.0),
		Bias:       true,
	})

	optimizer := training.NewAdam(0.001, 0.9, 0.999, 1e-8)
	trainer := training.NewBatchTrainer(optimizer, 1000, 200, 4)

	training, validation := trainingData.Split(0.75)
	trainer.Train(n, training, validation, 200000)

	for _, testCandidate := range testingData {
		fmt.Println(testCandidate.Input)
		fmt.Printf("final Test expected %v and was %v.\n",
			testCandidate.Response[0],
			n.Predict(testCandidate.Input)[0])
	}

	saveModelToFile(n)
}

func preprocess() (training.Examples, error) {
	f, err := os.Open(CSV_FILE)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvData := training.Examples{}

	csvReader := csv.NewReader(f)

	for i := 0; i < CSV_FILE_SKIP_FIRST_ROWS; i++ {
		_, err := csvReader.Read()
		if err != nil {
			return nil, err
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
		parsedLine, err := parseRecordToTrainingExample(rec)
		if err != nil {
			// we skip rows that are not parseable
			//fmt.Println("skipping row because: ", err.Error())
			continue
		}
		csvData = append(csvData, *parsedLine)
	}
	return csvData, nil
}

func parseRecordToTrainingExample(rec []string) (*training.Example, error) {
	housingMedianAge, err := strconv.ParseFloat(rec[2], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float housing median age: %s", rec[2])
	}
	totalRooms, err := strconv.ParseFloat(rec[3], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float total rooms: %s", rec[3])
	}
	totalBedrooms, err := strconv.ParseFloat(rec[4], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float total bedrooms: %s", rec[4])
	}
	population, err := strconv.ParseFloat(rec[5], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float population: %s", rec[5])
	}
	households, err := strconv.ParseFloat(rec[6], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float households: %s", rec[6])
	}
	medianIncome, err := strconv.ParseFloat(rec[7], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float median income: %s", rec[7])
	}
	var oceanProximity float64
	switch rec[9] {
	case "NEAR OCEAN":
		oceanProximity = 0
	case "NEAR BAY":
		oceanProximity = 1
	case "<1H OCEAN":
		oceanProximity = 2
	case "INLAND":
		oceanProximity = 3
	case "ISLAND":
		oceanProximity = 4
	default:
		return nil, fmt.Errorf("unknown ocean proximity: %s", rec[9])
	}
	medianHouseValue, err := strconv.ParseFloat(rec[8], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse float median house value: %s", rec[8])
	}
	houseValueOver140k := 0.0
	if medianHouseValue > 140000 {
		houseValueOver140k = 1
	}
	return &training.Example{
		Input: []float64{housingMedianAge,
			totalRooms, totalBedrooms, population,
			households, medianIncome, oceanProximity},
		Response: []float64{houseValueOver140k},
	}, nil
}

func saveModelToFile(network *deep.Neural) {
	marshalledNetwork, _ := network.Marshal()
	outputFile, err := os.OpenFile(
		fmt.Sprint(NETWORK_OUTPUT_FILE, time.Now().Format("_060102-1504")),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0600)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	if _, err = outputFile.WriteString(string(marshalledNetwork)); err != nil {
		panic(err)
	}
}
