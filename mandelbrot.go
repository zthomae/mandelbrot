package main

import (
	"image"
	"image/color"
	"image/gif"
	"os"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = iota
	blackIndex
)

// escapeIters returns the number of iterations
// needed to escape for the given c, or nIter if
// it doesn't escape within nIter iterations
func escapeIters(c complex128, nIter int) int {
	val := complex128(0)
	escape := float64(2)
	i := 0
	for i < nIter && real(val) < escape {
		val = val*val + c
		i++
	}
	return i
}

// mandelbrot draws an image of the mandelbrot set
// within the given real and imaginary bounds
func mandelbrot(center complex128, width, height float64) *image.Paletted {
	sX := 512 // TODO: Variable size
	sY := 512
	hX := 256.0
	hY := 256.0
	nIter := 1000
	lr := real(center) - width
	hr := real(center) + width
	li := imag(center) - height
	hi := imag(center) + height
	fromPos := func(x, y int) complex128 {
		r := real(center) + (float64(x)-hX)/float64(sX)*(hr-lr)
		i := imag(center) + (float64(y)-hY)/float64(sY)*(hi-li)
		var c complex128 = complex(r, i)
		return c
	}
	rect := image.Rect(0, 0, sX, sY)
	img := image.NewPaletted(rect, palette)
	for x := 0; x < sX; x++ {
		for y := 0; y < sY; y++ {
			c := fromPos(x, y)
			inSet := nIter == escapeIters(c, nIter)
			if inSet {
				img.SetColorIndex(x, y, blackIndex)
			} else {
				img.SetColorIndex(x, y, whiteIndex)
			}
		}
	}
	return img
}

type Frame struct {
	center complex128
	width  float64
	height float64
}

func scale(c complex128, s float64) complex128 {
	var res complex128 = complex(real(c)*s, imag(c)*s)
	return res
}

func main() {
	var startPos complex128 = complex(-1.02, 0.0)
	var endPos complex128 = complex(-1.31, 0.0)
	xs := [2]float64{0.25, 0.06}
	ys := [2]float64{0.25, 0.06}
	nFrames := 25
	frames := make([]Frame, nFrames)
	for i := 0; i < nFrames; i++ {
		frac := float64(i) / float64(nFrames-1)
		center := startPos + (scale(endPos, frac) - scale(startPos, frac))
		width := xs[0] + (xs[1]-xs[0])*frac
		height := ys[0] + (ys[1]-ys[0])*frac
		frames[i] = Frame{center, width, height}
	}
	anim := gif.GIF{LoopCount: 0}
	for _, vals := range frames {
		anim.Delay = append(anim.Delay, 8)
		anim.Image = append(anim.Image, mandelbrot(vals.center, vals.width, vals.height))
	}
	gif.EncodeAll(os.Stdout, &anim)
}
