package lib

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gobwas/glob/util/runes"
	"golang.org/x/exp/rand"
)

const (
	russianTelegraphAlpabet = "абвгдежзийклмнопрстуфхцчшщьыэюя "
	defaultStringLength     = 50
)

func GenerateStringByBlocks(path string) string {
	fmt.Println("123")
	text := russianTelegraphAlpabet
	if path != "" {
		text = GetDataFromFile(path)
	}

	symbolsCount := CountSymbols(text)

	siftedSymbolsCount := SiftRussianTelegraphAlpabet(symbolsCount)

	russianTotalCount := GetTotalCount(siftedSymbolsCount)

	symbolsFrequency := CalculateSymbolsFrequency(siftedSymbolsCount, russianTotalCount)

	return GenerateBlocks(text, symbolsFrequency)
}

func GenerateBlocks(text string, symbolsFrequency map[string]float64) string {
	oldBlock := PickNSymbolsInText(text, symbolsFrequency, 0, 1)
	fmt.Println("first bloc", oldBlock)
	for {
		newBlock := FindBlockInText(text, oldBlock, 2)
		fmt.Println(newBlock)
		if newBlock == "" {
			return oldBlock
		}

		oldBlock = newBlock
	}
}

func GenerateTwoString(path string, length int, shift int) (string, string) {
	text := russianTelegraphAlpabet
	if path != "" {
		text = GetDataFromFile(path)
	}

	symbolsCount := CountSymbols(text)

	siftedSymbolsCount := SiftRussianTelegraphAlpabet(symbolsCount)

	russianTotalCount := GetTotalCount(siftedSymbolsCount)

	symbolsFrequency := CalculateSymbolsFrequency(siftedSymbolsCount, russianTotalCount)

	entropyOutput := GenerateStringBasedOnFrequency(symbolsFrequency, length)

	oneShiftOutput := ""
	if text != russianTelegraphAlpabet {
		oneShiftOutput = GenerateStringWithNShift(text, symbolsFrequency, shift, length)
	}

	return entropyOutput, oneShiftOutput
}

func FindBlockInText(text string, block string, length int) string {
	index := runes.Index([]rune(text), []rune(block))
	fmt.Println("index", index)
	fmt.Println(string([]rune(text)[index : index+len(block)]))

	if index != -1 {
		return string([]rune(text)[index : index+length+len(block)])
	}

	return ""
}

func GenerateStringWithNShift(
	text string,
	symbolsFrequency map[string]float64,
	n int,
	length int,
) string {
	output := ""
	for i := 0; i < length; i++ {
		symbol := PickNSymbolInText(text, symbolsFrequency, n)
		output += symbol
	}

	return output
}

func PickNSymbolsInText(
	text string,
	symbolsFrequency map[string]float64,
	n int,
	length int,
) string {
	symbolsPositions := GetSymbolsPositions(text)

	symbol := PickSymbolBasedOnFrequency(symbolsFrequency)

	symbolPositions := symbolsPositions[symbol]

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Intn(len(symbolPositions))

	for {
		out := string([]rune(text)[symbolPositions[r]+n*2])
		if strings.Contains(russianTelegraphAlpabet, out) {
			return string([]rune(text)[symbolPositions[r]+n*2 : symbolPositions[r]+n*2+length*2])
		}
		n++
	}
}

func PickNSymbolInText(text string, symbolsFrequency map[string]float64, n int) string {
	symbolsPositions := GetSymbolsPositions(text)

	symbol := PickSymbolBasedOnFrequency(symbolsFrequency)

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

func GetSymbolsPositions(text string) map[string][]int {
	symbolsPositions := make(map[string][]int, 32)

	for position, symbol := range []rune(text) {
		lower := strings.ToLower(string(symbol))
		if strings.Contains(russianTelegraphAlpabet, lower) {
			symbolsPositions[lower] = append(symbolsPositions[lower], position)
		}
	}

	return symbolsPositions
}

func GenerateStringBasedOnFrequency(symbolsFrequency map[string]float64, length int) string {
	output := ""
	for i := 0; i < length; i++ {
		symbol := PickSymbolBasedOnFrequency(symbolsFrequency)
		output += symbol
	}

	return output
}

func CalculateEntropy(symbolsFrequency map[string]float64) float64 {
	entropy := 0.0
	for _, frequency := range symbolsFrequency {
		if frequency > 0 {
			entropy += frequency * math.Log2(frequency)
		}
	}
	return -entropy
}

func GetDataFromFile(path string) string {
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

func CountSymbols(text string) map[string]int {
	symbolsCount := make(map[string]int, 33)
	for _, symbol := range text {
		symbolsCount[strings.ToLower(string(symbol))]++
	}

	return symbolsCount
}

func SiftRussianTelegraphAlpabet(symbolsCount map[string]int) map[string]int {
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

func GetTotalCount(symbolsCount map[string]int) int {
	totalCount := 0
	for _, count := range symbolsCount {
		totalCount += count
	}
	return totalCount
}

func CalculateSymbolsFrequency(symbolsCount map[string]int, totalCount int) map[string]float64 {
	symbolsFrequency := make(map[string]float64, len(symbolsCount))
	for symbol, count := range symbolsCount {
		symbolsFrequency[symbol] = float64(count) / float64(totalCount)
	}
	return symbolsFrequency
}

func PickSymbolBasedOnFrequency(symbolsFrequency map[string]float64) string {
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
