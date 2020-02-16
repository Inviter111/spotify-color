package colorpicker

import (
	"image"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/pa-m/sklearn/cluster"
	"gonum.org/v1/gonum/mat"
)

// NULL variable
const NULL = ""

var lastImageURL string
var lastHex string

type imageMatrix [][]float64

// GetImageColor returns hex representation of best color
func GetImageColor(url string) string {
	if lastImageURL == url {
		return lastHex
	}

	img := getImage(url)
	lastImageURL = url

	resizedImg := imaging.Resize(img, 100, 100, imaging.Lanczos)
	mat := imageToMatrix(resizedImg)
	clt := cluster.KMeans{NClusters: 7}
	clt.Fit(mat, nil)
	centroids := clt.Centroids
	colors := createColors(centroids)
	colorfulness := make(map[string]float64)
	for _, color := range colors {
		colorfulness[color.toHex()] = color.colorfulness()
	}
	maxHex := maxColorfulness(colorfulness)
	lastHex = maxHex

	return maxHex
}

func getImage(url string) image.Image {
	resp, err := http.Get(url)
	if err != nil {
		printError(err)
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		printError(err)
	}

	return img
}

func imageToMatrix(i image.Image) imageMatrix {
	xSize, ySize := (i).Bounds().Size().X, (i).Bounds().Size().Y
	size := xSize * ySize
	mat := make([][]float64, size)
	for index := range mat {
		mat[index] = make([]float64, 3)
	}
	point := 0
	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			r, g, b, _ := (i).At(x, y).RGBA()
			mat[point][0] = float64(r)
			mat[point][1] = float64(g)
			mat[point][2] = float64(b)
			point++
		}
	}

	return mat
}

func printError(err error) {
	log.Fatalln("Error:", err)
	os.Exit(1)
}

func (i imageMatrix) Dims() (int, int) {
	return len(i), len(i[0])
}

func (i imageMatrix) At(x int, y int) float64 {
	return i[x][y]
}

func (i imageMatrix) T() mat.Matrix {
	return i
}

func maxColorfulness(m map[string]float64) string {
	uniqueElems := make([]float64, 0)
	for _, value := range m {
		if index := find(uniqueElems, value); index == -1 {
			uniqueElems = append(uniqueElems, value)
		}
	}
	if len(uniqueElems) == 1 {
		return NULL
	}
	var maxHex string
	max := float64(math.MinInt8)
	for hex, colorfulness := range m {
		if colorfulness > max {
			max = colorfulness
			maxHex = hex
		}
	}

	return maxHex
}

func find(a []float64, v float64) int {
	for i, n := range a {
		if v == n {
			return i
		}
	}

	return -1
}
