package main

import (
	"flag"
	"fmt"
	. "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
	"log"
	"os"
	"time"
)

const (
	scale              = 4
	near               = 1
	far                = 10
    defaultFovy        = 40.0
    defaultWidth       = 540
    defaultHeight      = 480
	defaultWorldColor  = "E0D9CC"
	defaultObjectColor = "F09500"
)

var (
	// render parameters
	eye         = V(3, 1, 1.75)
	center      = V(0, 0, 0)
	up          = V(0, 0, 1)
	worldColor  = HexColor(defaultWorldColor)  // TODO add command-line arg
	objectColor = HexColor(defaultObjectColor) // TODO add command-line arg
	// command-line args
	inputPath  = flag.String("input", "", "path to input file (required)")
	outputPath = flag.String("output", "", "path to output file (required)")
	width      = flag.Uint("width", defaultWidth, "output image width")
    height     = flag.Uint("height", defaultHeight, "output image height")
    fov        = flag.Float64("fov", defaultFovy, "field of view in degrees")
)

func main() {
	flag.Parse()
	if *inputPath == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Input file path required")
		os.Exit(0)
	}
	if *outputPath == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Output file path required")
		os.Exit(0)
	}

	// load mesh
	log.Println("Loading input file: " + *inputPath)
	mesh, err := LoadSTL(*inputPath)
	if err != nil {
		log.Fatalln("Failed to load input file: " + err.Error())
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()
	// rotate on z-axis
	mesh.Transform(Rotate(V(0, 0, 1.0), Radians(-90)))
	//mesh.SmoothNormalsThreshold(Radians(30))

	start := time.Now()
	context := NewContext(int(*width)*scale, int(*height)*scale)
	context.ClearColorBufferWith(worldColor)

	aspect := float64(*width) / float64(*height)
	matrix := LookAt(eye, center, up).Perspective(*fov, aspect, near, far)
	light := V(3.5, 1, 1.5).Normalize()

	// configure shader
	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = objectColor
	shader.DiffuseColor = Gray(0.8)
	shader.SpecularColor = Gray(0.55)
	shader.SpecularPower = 90
	context.Shader = shader

	log.Printf("Rendering image from mesh: %d triangles\n", len(mesh.Triangles))
    context.DrawMesh(mesh)
	log.Println("Render completed: " + time.Since(start).String())

	// save image
	image := context.Image()
	image = resize.Resize(*width, *height, image, resize.Bilinear)
    log.Printf("Writing image to output file: %s (%d x %d px)\n", *outputPath, *width, *height)
    err = SavePNG(*outputPath, image)
	if err != nil {
		log.Fatalln("Failed to save output image: " + err.Error())
	}
}
