package main

import (
	"fmt"
	"image/color"
)

// "fmt"
// "image"
// "log"
// "os"
//
// "image/color"
// _ "image/gif"
// _ "image/jpeg"
// "image/png"

import "github.com/gonum/matrix/mat64"

func main() {
	// imf, err := os.Open("magnitude.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer imf.Close()
	// fmt.Printf("reading %s...\n", imf.Name())
	//
	// // read image size
	// im, _, err := image.Decode(imf)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bounds := im.Bounds()
	// w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	//
	// // read image data
	// img := make([]color.Color, w*h)
	// for c := 0; c < w; c++ {
	// 	for r := 0; r < h; r++ {
	// 		img[c*h+r] = im.At(c, r)
	// 	}
	// }
	//
	// // process
	// fmt.Printf("processing %s[%dx%d]...", imf.Name(), w, h)
	// im_cone := image.NewGray(image.Rect(0, 0, w, h))
	// for c := 0; c < w; c++ {
	// 	for r := 0; r < h; r++ {
	// 		//var max float64
	// 		//var max_l int
	// 		//for l, _ := range imgs {
	// 		//	d := disp(imgs[l], c, r, width, height)
	// 		//	if d >= max {
	// 		//		max = d
	// 		//		max_l = l
	// 		//	}
	// 		//}
	// 		gray := float64(color.GrayModel.Convert(img[c*h+r]).(color.Gray).Y)
	// 		gray = gray * 1000000
	// 		im_cone.Set(c, r, color.Gray{uint8(gray)})
	// 	}
	// 	if c%10 == 0 {
	// 		fmt.Printf(".")
	// 	}
	// }
	//
	// fmt.Printf(" DONE\n")
	//
	// // save heights
	// cone, err := os.Create("magnitude_cone.png")
	// defer cone.Close()
	// png.Encode(cone, im_cone)
	// fmt.Printf("saving %s\n", cone.Name())

	//data := make([]float64, 9)
	//for i := range data {
	//	data[i] = rand.NormFloat64()
	//}
	//a := mat64.NewDense(3, 3, data)

	// Solve(a, b Matrix) error
	a := mat64.NewDense(3, 3, []float64{
		2, 0, 0,
		0, 1, 0,
		1, 0, 1,
	})
	b := mat64.NewVector(3, []float64{1, 2, 2})

	var x mat64.Vector
	if err := x.SolveVec(a, b); err != nil {
		fmt.Println("Matrix is near singular: ", err)
	}
	fmt.Println("Solve a * x = b")
	fmt.Printf("x = %0.4v\n", mat64.Formatted(&x, mat64.Prefix("    ")))

}

var area = 2

func disp(img []color.Color, col, row, width, height int) float64 {
	var (
		n     int
		x, x2 float64
	)
	for c := col - area; c < col+area; c++ {
		for r := row - area; r < row+area; r++ {
			if c < 0 || c >= width || r < 0 || r >= height || (col-c)*(col-c)+(row-r)*(row-r) > area*area {
				// continue
			} else {
				n += 1
				gray := float64(color.GrayModel.Convert(img[c*height+r]).(color.Gray).Y)
				x += gray
				x2 += gray * gray
			}
		}
	}
	return float64(n)*x2/(x*x) - 1
}
