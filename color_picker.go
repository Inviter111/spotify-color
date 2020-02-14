package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pa-m/sklearn/cluster"
	"gonum.org/v1/gonum/mat"
)

type imageMatrix [][]float64

func main() {
	start := time.Now()
	img := getImage()
	// fmt.Println(img.At(63, 64).RGBA())
	mat := imageToMatrix(&img)
	// fmt.Println(mat)
	clt := cluster.KMeans{NClusters: 3}
	clt.Fit(mat, nil)
	centroids := clt.Centroids
	for i := 0; i < 3; i++ {
		r, g, b := int(centroids.At(i, 0)), int(centroids.At(i, 1)), int(centroids.At(i, 2))
		rHex, gHex, bHex := fmt.Sprintf("%02X", r/0x101), fmt.Sprintf("%02X", g/0x101), fmt.Sprintf("%02X", b/0x101)
		fmt.Printf("#%s%s%s\n", rHex, gHex, bHex)
	}
	elapsed := time.Since(start)
	fmt.Println("Execution time:", elapsed.Seconds())
}

func getImage() image.Image {
	resp, err := http.Get("https://i.scdn.co/image/ab67616d00001e02466f56d5f68eec9b0866e894")
	if err != nil {
		printError(err)
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		printError(err)
	}

	return img
}

func imageToMatrix(i *image.Image) imageMatrix {
	xSize, ySize := (*i).Bounds().Size().X, (*i).Bounds().Size().Y
	size := xSize * ySize
	mat := make([][]float64, size)
	for index := range mat {
		mat[index] = make([]float64, 3)
	}
	point := 0
	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			r, g, b, _ := (*i).At(x, y).RGBA()
			mat[point][0] = float64(r)
			mat[point][1] = float64(g)
			mat[point][2] = float64(b)
			point++
		}
	}

	return mat
}

// func bestColor() {
// 	image := parseImage()
// 	clt := &cluster.KMeans{}
// 	// clt.Fit()
// }

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
