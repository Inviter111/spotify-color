package colorpicker

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

type hexColor struct {
	r int
	g int
	b int
}

func createColors(d *mat.Dense) []hexColor {
	size, _ := d.Dims()
	colors := make([]hexColor, size)
	for i := 0; i < size; i++ {
		colors[i].r, colors[i].g, colors[i].b = int(d.At(i, 0)), int(d.At(i, 1)), int(d.At(i, 2))
	}

	return colors
}

func (c hexColor) toHex() string {
	r, g, b := c.r, c.g, c.b
	rHex, gHex, bHex := fmt.Sprintf("%02X", r/0x101), fmt.Sprintf("%02X", g/0x101), fmt.Sprintf("%02X", b/0x101)
	return fmt.Sprintf("%s%s%s", rHex, gHex, bHex)
}

func (c hexColor) colorfulness() float64 {
	r, g, b := float64(c.r), float64(c.g), float64(c.b)

	rg := math.Abs(r - g)
	yb := math.Abs(0.5*(r+g) - b)

	root := math.Sqrt(math.Pow(rg, 2) + math.Pow(yb, 2))

	return 0.3 * root
}
