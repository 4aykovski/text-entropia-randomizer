package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/4aykovski/text-entropia-randomizer/lib"
	"github.com/4aykovski/text-entropia-randomizer/ui"
	"golang.org/x/exp/rand"
)

const (
	russianTelegraphAlpabet = "абвгдежзийклмнопрстуфхцчшщьыэюя "
	defaultStringLength     = 50
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		ui.Run()
	} else {
		text := russianTelegraphAlpabet

		path := flag.String("f", "", "path to file")
		length := flag.Int("l", defaultStringLength, "length of generated string")
		flag.Parse()

		if *path != "" {
			text = getDataFromFile(*path)
		}

		symbolsCount := countSymbols(text)

		siftedSymbolsCount := siftRussianTelegraphAlpabet(symbolsCount)

		russianTotalCount := getTotalCount(siftedSymbolsCount)

		symbolsFrequency := calculateSymbolsFrequency(siftedSymbolsCount, russianTotalCount)

		entropyOutput := generateStringBasedOnFrequency(symbolsFrequency, *length)
		oneShiftOutput := generateStringWithNShift(text, symbolsFrequency, 1, *length)
		twoShiftOutput := generateStringWithNShift(text, symbolsFrequency, 2, *length)

		block := lib.GenerateStringByBlocks(*path)
		fmt.Println(block)

		fmt.Println("entropy output \n\t", entropyOutput)
		fmt.Println("one shift output \n\t ", oneShiftOutput)
		fmt.Println("two shift output \n\t", twoShiftOutput)
	}
}

func generateStringWithNShift(
	text string,
	symbolsFrequency map[string]float64,
	n int,
	length int,
) string {
	output := ""
	for i := 0; i < length; i++ {
		symbol := pickNSymbolInText(text, symbolsFrequency, n)
		output += symbol
	}

	return output
}

func pickNSymbolInText(text string, symbolsFrequency map[string]float64, n int) string {
	symbolsPositions := getSymbolsPositions(text)

	symbol := pickSymbolBasedOnFrequency(symbolsFrequency)

	symbolPositions := symbolsPositions[symbol]

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Intn(len(symbolPositions))

	for {
		out := string([]rune(text)[symbolPositions[r]+n*2])
		if strings.Contains(russianTelegraphAlpabet, out) {
			return out
		}
		n++
	}
}

func getSymbolsPositions(text string) map[string][]int {
	symbolsPositions := make(map[string][]int, 32)

	for position, symbol := range []rune(text) {
		lower := strings.ToLower(string(symbol))
		if strings.Contains(russianTelegraphAlpabet, lower) {
			symbolsPositions[lower] = append(symbolsPositions[lower], position)
		}
	}

	return symbolsPositions
}

func generateStringBasedOnFrequency(symbolsFrequency map[string]float64, length int) string {
	output := ""
	for i := 0; i < length; i++ {
		symbol := pickSymbolBasedOnFrequency(symbolsFrequency)
		output += symbol
	}

	return output
}

func calculateEntropy(symbolsFrequency map[string]float64) float64 {
	entropy := 0.0
	for _, frequency := range symbolsFrequency {
		if frequency > 0 {
			entropy += frequency * math.Log2(frequency)
		}
	}
	return -entropy
}

func getDataFromFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var text string
	data := make([]byte, 1024)
	for {
		n, err := file.Read(data)
		if err != nil {
			break
		}
		text += string(data[:n])
	}

	return text
}

func countSymbols(text string) map[string]int {
	symbolsCount := make(map[string]int, 33)
	for _, symbol := range text {
		symbolsCount[strings.ToLower(string(symbol))]++
	}

	return symbolsCount
}

func siftRussianTelegraphAlpabet(symbolsCount map[string]int) map[string]int {
	siftedSymbolsCount := make(map[string]int, 32)

	for symbol, count := range symbolsCount {
		if symbol == "ё" {
			currentCouunt := siftedSymbolsCount["е"]
			siftedSymbolsCount["е"] = currentCouunt + count
			continue
		}

		if symbol == "ъ" || symbol == "ь" {
			currentCouunt := siftedSymbolsCount["ь"]
			siftedSymbolsCount["ь"] = currentCouunt + count
			continue
		}

		if strings.Contains(russianTelegraphAlpabet, symbol) {
			siftedSymbolsCount[symbol] = count
		}
	}

	return siftedSymbolsCount
}

func getTotalCount(symbolsCount map[string]int) int {
	totalCount := 0
	for _, count := range symbolsCount {
		totalCount += count
	}
	return totalCount
}

func calculateSymbolsFrequency(symbolsCount map[string]int, totalCount int) map[string]float64 {
	symbolsFrequency := make(map[string]float64, len(symbolsCount))
	for symbol, count := range symbolsCount {
		symbolsFrequency[symbol] = float64(count) / float64(totalCount)
	}
	return symbolsFrequency
}

func pickSymbolBasedOnFrequency(symbolsFrequency map[string]float64) string {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Float64()
	currentWeight := 0.0

	for symbol, frequency := range symbolsFrequency {
		currentWeight += frequency
		if r < currentWeight {
			return symbol
		}
	}

	return ""
}
