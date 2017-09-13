package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
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
	imf, err := os.Open("magnitude_log.png")
	if err != nil {
		log.Fatal(err)
	}
	defer imf.Close()
	fmt.Printf("reading %s...\n", imf.Name())

	// read image size
	im, _, err := image.Decode(imf)
	if err != nil {
		log.Fatal(err)
	}
	bounds := im.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	// (ax^2+by^2+cxy+d-f(x,y))*x^2 = 0
	//                          y^2 = 0
	//                          xy  = 0
	//                          1   = 0
	//
	// x4   x2y2 x3y  x2  | fx2
	// x2y2 y4   xy3  y2  | fy2
	// x3y  xy3  x2y2 xy  | fxy
	// x2   y2   xy   1   | f

	x2 := float64(0)
	y2 := float64(0)
	x4 := float64(0)
	y4 := float64(0)
	xy := float64(0)
	x3y := float64(0)
	xy3 := float64(0)
	x2y2 := float64(0)

	f := float64(0)
	fx2 := float64(0)
	fy2 := float64(0)
	fxy := float64(0)

	n := 0
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			x := float64(c - w/2)
			y := float64(r - h/2)
			if x*x+y*y >= float64(w*h)/100.0 {
				n = n + 1
			}
		}
	}

	dd := 1 / float64(n)
	fmt.Printf("n: %d\n", n)
	fmt.Printf("dd: %f\n", dd)

	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			f_im := float64(color.GrayModel.Convert(im.At(c, r)).(color.Gray).Y)

			x := float64(c - w/2)
			y := float64(r - h/2)

			//if x*x+y*y < float64((w*h))/36.0 && x*x+y*y >= float64((w+h)*(w+h))/360.0 {
			if x*x+y*y >= float64(w*h)/100.0 {
				x2 += x * x * dd
				y2 += y * y * dd
				x4 += x * x * x * x * dd
				y4 += y * y * y * y * dd
				xy += x * y * dd
				x3y += x * x * x * y * dd
				xy3 += x * y * y * y * dd
				x2y2 += x * x * y * y * dd

				f += f_im * dd
				fx2 += f_im * x * x * dd
				fy2 += f_im * y * y * dd
				fxy += f_im * x * y * dd
			}
		}
	}

	fmt.Printf("xy: %f\n", xy)

	a := mat64.NewDense(4, 4, []float64{
		x4, x2y2, x3y, x2,
		x2y2, y4, xy3, y2,
		x3y, xy3, x2y2, xy,
		x2, y2, xy, 1,
	})
	b := mat64.NewVector(4, []float64{fx2, fy2, fxy, f})

	var x mat64.Vector
	if err := x.SolveVec(a, b); err != nil {
		fmt.Println("Matrix is near singular: ", err)
	}
	fmt.Println("Solve a * x = b")
	fmt.Printf("x = %0.4v\n", mat64.Formatted(&x, mat64.Prefix("    ")))

	fmt.Printf("processing %s[%dx%d]...", imf.Name(), w, h)
	im_cone := image.NewGray(image.Rect(0, 0, w, h))
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			//var max float64
			//var max_l int
			//for l, _ := range imgs {
			//	d := disp(imgs[l], c, r, width, height)
			//	if d >= max {
			//		max = d
			//		max_l = l
			//	}
			//}
			// gray := float64(color.GrayModel.Convert(im.At(c, r)).(color.Gray).Y)
			// gray = gray * 1000000

			xa := x.At(0, 0)
			xb := x.At(1, 0)
			xc := x.At(2, 0)
			xd := x.At(3, 0)

			x := float64(c - w/2)
			y := float64(r - h/2)

			gray := float64(0)
			if x*x+y*y >= float64(w*h)/100.0 {
				gray = xa*x*x + xb*y*y + xc*x*y + xd
			}

			im_cone.Set(c, r, color.Gray{uint8(gray)})
		}
		if c%10 == 0 {
			fmt.Printf(".")
		}
	}

	fmt.Printf(" DONE\n")

	// save heights
	cone, err := os.Create("magnitude_cone.png")
	defer cone.Close()
	png.Encode(cone, im_cone)
	fmt.Printf("saving %s\n", cone.Name())
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
