package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
	"strconv"
	"sync"
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
func mandelbrot(sX, sY, nIter int,
	center complex128,
	width, height float64) *image.Paletted {

	hX := float64(sX) / 2.0
	hY := float64(sY) / 2.0
	lr := real(center) - width/2.0
	hr := real(center) + width/2.0
	li := imag(center) - height/2.0
	hi := imag(center) + height/2.0
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

// State stores all of the data specifying an animation
type State struct {
	startPos, endPos              complex128
	startZoom, endZoom            [2]float64
	filename                      string
	sX, sY, nIter, nFrames, delay int
}

// parse command-line arguments, returning a valid State struct
func args() State {
	s := State{}
	startPosSet := false
	endPosSet := false
	startZoomSet := false
	endZoomSet := false
	sizeSet := false
	nIterSet := false
	nFramesSet := false
	delaySet := false
	argc := len(os.Args)
	i := 1
	for i < argc {
		var err error
		switch os.Args[i] {
		case "--startPos":
			var r, im float64 // initialized in inner scopes
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument(s) to --startPos")
				os.Exit(1)
			}
			r, err := strconv.ParseFloat(os.Args[i+1], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --startPos\n", os.Args[i+1])
				os.Exit(1)
			}
			if i == argc-2 {
				im = 0.0
				i += 2
			} else {
				im, err = strconv.ParseFloat(os.Args[i+2], 64)
				if os.Args[i+2][0] == '-' {
					im = 0.0
					i += 2
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --startPos\n", os.Args[i+2])
					os.Exit(1)
				} else {
					i += 3
				}
			}
			s.startPos = complex(r, im)
			startPosSet = true
		case "--endPos":
			var r, im float64 // initialized in inner scopes
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument(s) to --endPos")
				os.Exit(1)
			}
			r, err := strconv.ParseFloat(os.Args[i+1], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --endPos\n", os.Args[i+1])
				os.Exit(1)
			}
			if i == argc-2 {
				im = 0.0
				i += 2
			} else {
				im, err = strconv.ParseFloat(os.Args[i+2], 64)
				if os.Args[i+2][0] == '-' {
					im = 0.0
					i += 2
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --endPos\n", os.Args[i+2])
					os.Exit(1)
				} else {
					i += 3
				}
			}
			s.endPos = complex(r, im)
			endPosSet = true
		case "--startZoom":
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument(s) to --startZoom")
				os.Exit(1)
			}
			s.startZoom[0], err = strconv.ParseFloat(os.Args[i+1], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --startZoom\n", os.Args[i+1])
				os.Exit(1)
			}
			if i == argc-2 {
				s.startZoom[1] = s.startZoom[0]
				i += 2
			} else {
				s.startZoom[1], err = strconv.ParseFloat(os.Args[i+2], 64)
				if os.Args[i+2][0] == '-' {
					s.startZoom[1] = s.startZoom[0]
					i += 2
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --startZoom\n", os.Args[i+2])
					os.Exit(1)
				} else {
					i += 3
				}
			}
			startZoomSet = true
		case "--endZoom":
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument(s) to --endZoom")
				os.Exit(1)
			}
			s.endZoom[0], err = strconv.ParseFloat(os.Args[i+1], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --endZoom\n", os.Args[i+1])
				os.Exit(1)
			}
			if i == argc-2 {
				s.endZoom[1] = s.endZoom[0]
				i += 2
			} else {
				s.endZoom[1], err = strconv.ParseFloat(os.Args[i+2], 64)
				if os.Args[i+2][0] == '-' {
					s.endZoom[1] = s.endZoom[0]
					i += 2
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --endZoom\n", os.Args[i+2])
					os.Exit(1)
				} else {
					i += 3
				}
			}
			endZoomSet = true
		case "--output":
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument to --output")
				os.Exit(1)
			}
			s.filename = os.Args[i+1]
			i += 2
		case "--size":
			var sX, sY int64 // initialized in inner scopes
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument(s) to --size")
				os.Exit(1)
			}
			sX, err = strconv.ParseInt(os.Args[i+1], 10, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --size\n", os.Args[i+1])
			}
			if i == argc-2 {
				sY = sX
				i += 2
			} else {
				sY, err = strconv.ParseInt(os.Args[i+2], 10, 0)
				if os.Args[i+2][0] == '-' {
					sY = sX
					i += 2
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --size\n", os.Args[i+2])
					os.Exit(1)
				} else {
					i += 3
				}
			}
			s.sX = int(sX)
			s.sY = int(sY)
			sizeSet = true
		case "--iters":
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument to --iters")
				os.Exit(1)
			}
			nIter, err := strconv.ParseInt(os.Args[i+1], 10, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --iters\n", os.Args[i+1])
				os.Exit(1)
			}
			s.nIter = int(nIter)
			i += 2
			nIterSet = true
		case "--frames":
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument to --frames")
				os.Exit(1)
			}
			nFrames, err := strconv.ParseInt(os.Args[i+1], 10, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --frames\n", os.Args[i+1])
				os.Exit(1)
			}
			if nFrames < 1 {
				fmt.Fprintln(os.Stderr, "number of frames must be at least 1")
				os.Exit(1)
			}
			s.nFrames = int(nFrames)
			i += 2
			nFramesSet = true
		case "--delay":
			if i == argc-1 {
				fmt.Fprintln(os.Stderr, "expected argument to --delay")
				os.Exit(1)
			}
			delay, err := strconv.ParseInt(os.Args[i+1], 10, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot parse %s as argument to --delay\n", os.Args[i+1])
				os.Exit(1)
			}
			if delay < 1 {
				fmt.Fprintf(os.Stderr, "delay time must be at least 1")
				os.Exit(1)
			}
			s.delay = int(delay)
			i += 2
			delaySet = true
		case "--test":
			// TODO: This is not a perfect animation...
			s.startPos = complex(-1.0, 0.0)
			s.endPos = complex(-1.31, 0.0)
			s.startZoom[0] = 0.5
			s.startZoom[1] = 0.5
			s.endZoom[0] = 0.12
			s.endZoom[1] = 0.12
			s.sX = 512
			s.sY = 512
			s.nIter = 1000
			s.nFrames = 25
			s.delay = 8
			s.filename = "test.gif" // good decision?
			return s
		default:
			fmt.Fprintf(os.Stderr, "unexpected argument %s\n", os.Args[i])
			os.Exit(1)
		}
	}
	if !startPosSet {
		fmt.Fprintln(os.Stderr, "need to specify start position")
		os.Exit(1)
	} else if !endPosSet {
		s.endPos = s.startPos
	}
	if !startZoomSet && !endZoomSet {
		fmt.Fprintln(os.Stderr, "need to give start zoom")
		os.Exit(1)
	} else if !endZoomSet {
		s.endZoom[0] = s.startZoom[0]
		s.endZoom[1] = s.startZoom[1]
	}
	if !sizeSet {
		s.sX = 512
		s.sY = 512
	}
	if !nIterSet {
		s.nIter = 1000
	}
	if !nFramesSet && !endZoomSet {
		s.nFrames = 1
	} else if !nFramesSet {
		s.nFrames = 25
	}
	if s.nFrames > 1 && !endZoomSet && !endPosSet {
		fmt.Fprintln(os.Stderr, "setting frames argument to 1 due to lack of movement")
		s.nFrames = 1
	}
	if s.nFrames == 1 {
		if endZoomSet {
			fmt.Fprintln(os.Stderr, "frames set to 1; ignoring end zoom")
		}
		if endPosSet {
			fmt.Fprintln(os.Stderr, "frames set to 1; ignoring end position")
		}
	}
	if !delaySet {
		s.delay = 8
	}
	return s
}

// Frame holds the parameters for a given image (position and zoom)
type Frame struct {
	center complex128
	width  float64
	height float64
}

// scale a complex number by a scalar float
func scale(c complex128, s float64) complex128 {
	var res complex128 = complex(real(c)*s, imag(c)*s)
	return res
}

// create animation according to the State struct
func (s State) animate() {
	xs := [2]float64{s.startZoom[0], s.endZoom[0]}
	ys := [2]float64{s.startZoom[1], s.endZoom[1]}
	frames := make([]Frame, s.nFrames)
	for i := 0; i < s.nFrames; i++ {
		denom := s.nFrames
		if denom > 1 {
			denom--
		}
		frac := float64(i) / float64(denom)
		center := s.startPos + (scale(s.endPos, frac) - scale(s.startPos, frac))
		width := xs[0] + (xs[1]-xs[0])*frac
		height := ys[0] + (ys[1]-ys[0])*frac
		frames[i] = Frame{center, width, height}
	}

	// spawn goroutines to draw the frames, waiting for them all to finish
	anim := gif.GIF{LoopCount: 0}
	anim.Delay = make([]int, s.nFrames)
	anim.Image = make([]*image.Paletted, s.nFrames)
	var wg sync.WaitGroup
	wg.Add(s.nFrames)
	for i, vals := range frames {
		anim.Delay[i] = s.delay
		go func(i int, vals Frame) {
			defer wg.Done()
			anim.Image[i] = mandelbrot(s.sX, s.sY, s.nIter,
				vals.center, vals.width, vals.height)
		}(i, vals)
	}
	wg.Wait()

	var out io.Writer
	if s.filename == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(s.filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		out = bufio.NewWriter(f)
	}
	gif.EncodeAll(out, &anim)
}

func main() {
	args().animate()
}
