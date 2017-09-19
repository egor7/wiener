package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"image/color"
	"image/png"
)

// "fmt"
// "image"
// "log"
// "os"
//
// "image/color"
// _ "image/gif"
// _ "image/jpeg"

import "github.com/gonum/matrix/mat64"

func main() {

	// u
	f_u, err := os.Open("u.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f_u.Close()
	im_u, _, err := image.Decode(f_u)
	if err != nil {
		log.Fatal(err)
	}
	bounds := im_u.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	h_u := make([]float64, w*h)
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			h_u[c*h+r] = float64(color.GrayModel.Convert(im_u.At(c, r)).(color.Gray).Y)
		}
	}
	fmt.Printf("read %s, size [%dx%d]\n", f_u.Name(), w, h)

	// n
	h_n := make([]float64, w*h)
	a := 2.0
	a2 := a * a
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			x := float64(c - w/2)
			y := float64(r - h/2)
			x2 := x * x
			y2 := y * y
			h_n[c*h+r] = math.Exp((-x2-y2)/a2) / (a2 * math.Pi)
		}
	}

	// s
	h_s := make([]float64, w*h)
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			h_s[c*h+r] = smear(h_u, h_n, w, h, c, r)
		}
		if c%10 == 0 {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")

	// s
	im_s := image.NewGray(image.Rect(0, 0, w, h))
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			h := h_s[c*h+r]
			im_s.Set(c, r, color.Gray{uint8(h)})
		}
	}
	f_s, err := os.Create("s.png")
	defer f_s.Close()
	png.Encode(f_s, im_s)
	fmt.Printf("saved %s\n", f_s.Name())

	// n
	im_n := image.NewGray(image.Rect(0, 0, w, h))
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			h := h_n[c*h+r]
			im_n.Set(c, r, color.Gray{uint8(h)})
		}
	}
	f_n, err := os.Create("n.png")
	defer f_n.Close()
	png.Encode(f_n, im_n)
	fmt.Printf("saved %s\n", f_n.Name())

	// ~u
	im_tu := image.NewGray(image.Rect(0, 0, w, h))
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			h := h_u[c*h+r]
			im_tu.Set(c, r, color.Gray{uint8(h)})
		}
	}
	f_tu, err := os.Create("tu.png")
	defer f_tu.Close()
	png.Encode(f_tu, im_tu)
	fmt.Printf("saved %s\n", f_tu.Name())
}

func smear(orig, noise []float64, w, h int, oc, or int) float64 {
	s := 0.0
	for c := oc - 5; c < oc+5; c++ {
		for r := or - 5; r < or+5; r++ {
			nc := c - oc + w/2
			nr := r - or + h/2

			// tile
			cc := (c + w) % w
			rr := (r + h) % h
			s += orig[cc*h+rr] * noise[nc*h+nr]
		}
	}
	return s
}

var area = 2

// side = 1 => exterior
// side = -1 => interior
func approx(img []float64, w int, h int, dr int, side int) (xa, xb, xc, delta float64) {
	xa = 0
	xb = 0
	xc = 0
	delta = 0
	// (ax^2+by^2+cxy+d-f(x,y))*x^2 = 0
	//                          y^2 = 0
	//                          xy  = 0
	//                          1   = 0
	//
	// x4   x2y2 x3y  x2  | fx2
	// x2y2 y4   xy3  y2  | fy2
	// x3y  xy3  x2y2 xy  | fxy
	// x2   y2   xy   1   | f

	// instead use a plate:

	// (ax + by + c - f(x,y))^2 ~> min
	//
	// (ax + by + c - f(x,y))x
	// (ax + by + c - f(x,y))y
	// (ax + by + c - f(x,y))
	//
	// x2 xy x1 | fx
	// xy y2 y1 | fy
	// x1 y1  1 | f

	n := 0
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			x := c - w/2
			y := r - h/2
			if x >= 0 && y >= 0 && side*(x*x+y*y) > side*(dr*dr) {
				n = n + 1
			}
		}
	}

	var (
		x2, y2, xy, x1, y1 float64
		f, fx, fy          float64
	)
	dd := 1 / float64(n)
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			f_im := img[c*h+r]

			x := c - w/2
			y := r - h/2

			if x >= 0 && y >= 0 && side*(x*x+y*y) > side*(dr*dr) {
				x2 += float64(x*x) * dd
				y2 += float64(y*y) * dd
				xy += float64(x*y) * dd
				x1 += float64(x) * dd
				y1 += float64(y) * dd

				f += f_im * dd
				fx += float64(x) * f_im * dd
				fy += float64(y) * f_im * dd
			}
		}
	}

	a := mat64.NewDense(3, 3, []float64{
		x2, xy, x1,
		xy, y2, y1,
		x1, y1, 1,
	})
	b := mat64.NewVector(3, []float64{fx, fy, f})

	var x mat64.Vector
	if err := x.SolveVec(a, b); err != nil {
		//fmt.Println("Matrix is near singular: ", err)
	}
	//fmt.Printf("x = %0.4v\n", mat64.Formatted(&x, mat64.Prefix("    ")))

	xa = x.At(0, 0)
	xb = x.At(1, 0)
	xc = x.At(2, 0)
	for c := 0; c < w; c++ {
		for r := 0; r < h; r++ {
			x := float64(c - w/2)
			y := float64(r - h/2)

			if x >= 0 && y >= 0 && float64(side)*(x*x+y*y) > float64(side*dr*dr) {
				delta += math.Sqrt((xa*x+xb*y+xc-img[c*h+r])*(xa*x+xb*y+xc-img[c*h+r])) * dd
			}
		}
	}

	return
}
