# mandelbrot

A small program for generating gifs of the Mandelbrot set. Has command-line
options for most of the parameters describing an animation:

* Start/end position
* Start/end zoom
* Image width/height
* Number of iterations
* Number of frames
* Time delay between frames
* Output filename (defaults to standard out)

A test animation can be created with the `--test` option.

## USAGE

The following flags are valid:

* `--startPos <real> [<imag>]`: Set the center position at the start
      of the animation. If only `real` given, set `imag` to `real`. Mandatory.
* `--endPos <real> [<imag>]`: Set the center position at the end of
      the animation. If only `real` given, set `imag` to `real`. If not given,
	  this will be set to the starting position.
* `--startZoom <real> [<imag>]`: Set the zoom (the width/height visible
      in the image) at the start of the animation. If `imag` is not given,
	  it will be set to 0. Mandatory
* `--endZoom <real> [<imag>]`: Set the zoom at the end of the animation. If
      `imag` is not given, it will be set to 0. If not given, this will be
	  set to the starting zoom.
* `--size <x> [<y>]`: Set the size (in pixels) of the image drawn. If `y` is
      not given, it will be set to `x`. Defaults to 512x512.
* `--iters <n>`: Set the number of iterations to use when testing for
      membership in the Mandelbrot set. Defaults to 1000.
* `--frames <n>`: Set the number of frames that will be in the animation.
      Defaults to 25.
* `--delay <n>`: Sets the delay between frames (in 100ths of a second).
      Defaults to 8
* `--output <f>`: Sets the output file to `f`. If not given, will output to
      standard out.
* `--help`: Displays this message.
* `--test`: Runs the program with test parameters

Usage information can be printed by either invoking the command with no flags
or by invoking with `--help`.

## TODO:

* Expand the color palette and give other drawing options
* Testing?

## LICENSE

MIT
