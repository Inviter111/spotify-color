package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

type imageMatrix [][]float64

func main() {
	img := getImage()
	// fmt.Println(img.At(63, 64).RGBA())
	mat := imageToMatrix(&img)
	fmt.Println(mat)
}

func getImage() image.Image {
	resp, err := http.Get("https://i.scdn.co/image/5a73a056d0af707b4119a883d87285feda543fbb")
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
	// fmt.Println((*i).At(0, 0).RGBA())
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

// func (i imageMatrix) Dims() (int, int) {
// 	return i.Bounds().Size().X, i.Bounds().Size().Y
// }

// func (i imageMatrix) At(i, j int) float64 {
// 	return image.Image(i).At(i, j)
// }
