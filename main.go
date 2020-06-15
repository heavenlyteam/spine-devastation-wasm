// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"encoding/json"
	"math"
	"strings"
	"syscall/js"
)

type Pixel struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

var c chan bool

type jsInput [][]Pixel

const LR = 682.5
const LG = 532.5
const LB = 460

func addFunction(this js.Value, p []js.Value) interface{} {
	var (
		payload jsInput
		err     error
	)

	input := p[0].String()
	if err = json.NewDecoder(strings.NewReader(input)).Decode(&payload); err != nil {
		panic("Cant decode")
	}

	result := make([][]float64, len(payload))

	for i, row := range payload {
		result[i] = make([]float64, cap(row))
		for j, pixel := range row {
			r := math.Pow(pixel.R/LR, 2)
			g := math.Pow(pixel.G/LG, 2)
			b := math.Pow(pixel.B/LB, 2)

			pixelValue := math.Sqrt(r+g+b) / 255

			r = math.Pow(1/LR, 2)
			g = math.Pow(1/LG, 2)
			b = math.Pow(1/LB, 2)
			pixelValue *= 1 / math.Sqrt(r+g+b)

			result[i][j] = pixelValue
		}
	}

	var respString []byte
	if respString, err = json.Marshal(result); err != nil {
		panic("Cant  encode response")
	}

	return js.ValueOf(string(respString))
}

func main() {
	c = make(chan bool)
	js.Global().Set("pixelCalc", js.FuncOf(addFunction))

	<-c
}
