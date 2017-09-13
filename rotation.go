package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/andlabs/ui"
	"github.com/davecgh/go-spew/spew"
	flag "github.com/ogier/pflag"
)

var err error
var R *Rotation

// A point is a point in three dimensional space
type point struct {
	x float64
	y float64
	z float64
}

// A plane is an ordered set of points, each of which in an index to the point collection. The points are ordered such
// that p1->p2 forms vector A and points p2->p3 form vector B. With the points being ordered it is possible to calculate
// a vector C that is the othogonal normal to the plane.
type plane struct {
	p1 int
	p2 int
	p3 int
}

type Config struct {
	debug     int
	filename  string
	wireframe bool
	delay     float64
}

type Rotation struct {
	pointCount  int
	planeCount  int
	points      map[int]point
	planes      map[int]plane
	config      Config
}

func NewRotation() *Rotation {
	o := Rotation{
		0,
		0,
		make(map[int]point),
		make(map[int]plane),
		Config{
			0,
			"",
			false,
			0.0,
		},
	}
	return &o
}

type Rotatable interface {
	ReadOpts()
}

func (r *Rotation) ReadOpts() error {
	flag.StringVarP(&r.config.filename, "filename", "f", "", "Data file to load")
	flag.Float64VarP(&r.config.delay, "delay", "d", 0.0, "Delay to add to rotation. (Default: 0) ")
	flag.BoolVarP(&r.config.wireframe, "wireframe", "w", false, "Show object as a wirefram? (Default: no)")
	flag.IntVarP(&r.config.debug, "debug", "g", 0, "Debug logging level. 0=none, 1=some (Default: 0)")

	flag.Parse()

	if r.config.filename == "" {
		return errors.New("Missing filename")
	}

	if r.config.delay < 0 {
		return errors.New("Delay must be positive.")
	}
	return nil

}

func (r *Rotation) ReadFile() error {

	f, e := os.Open(r.config.filename)
	if e != nil {
		panic(e)
	}
	fs := bufio.NewScanner(f)

	fs.Scan()
	r.pointCount, _ = strconv.Atoi(fs.Text())
	if r.config.debug > 0 {
		log.Println(r.pointCount, " points to be read.")
	}

	for i := 0; i < r.pointCount; i++ {
		if r.config.debug > 0 {
			log.Println("reading points itteration: ", i)
		}
		line := ""
		for line == "" {
			fs.Scan()
			line = strings.Trim(fs.Text(), " ")
		}
		parts := strings.Fields(line)

		x, e1 := strconv.ParseFloat(parts[0], 64)
		y, e2 := strconv.ParseFloat(parts[1], 64)
		z, e3 := strconv.ParseFloat(parts[2], 64)

		if r.config.debug > 1 {
			spew.Dump(x, y, z, e1, e2, e3)
		}
		r.points[i] = point{
			x,
			y,
			z,
		}
	}

	fs.Scan()
	r.planeCount, _ = strconv.Atoi(fs.Text())
	if r.config.debug > 0 {
		log.Println(r.planeCount, " planes to be read.")
	}

	for i := 0; i < r.planeCount; i++ {
		if r.config.debug > 0 {
			log.Println("reading planes itteration: ", i)
		}
		line := ""
		for line == "" {
			fs.Scan()
			line = strings.Trim(fs.Text(), " ")
		}
		parts := strings.Fields(line)

		v1, e1 := strconv.Atoi(parts[0])
		v2, e2 := strconv.Atoi(parts[1])
		v3, e3 := strconv.Atoi(parts[2])

		if r.config.debug > 1 {
			spew.Dump(v1, v2, v3, e1, e2, e3)
		}
		r.planes[i] = plane{
			v1,
			v2,
			v3,
		}
	}

	if r.config.debug > 0 {
		var keys []int

		for i := range r.points {
			keys = append(keys, i)
		}
		sort.Ints(keys)
		log.Println("Points:")
		for _, i := range keys {
			log.Printf("i: %2d x: %3d y: %3d z: %3d \n", i, r.points[i].x, r.points[i].y, r.points[i].z)
		}

		for i := range r.planes {
			keys = append(keys, i)
		}
		sort.Ints(keys)
		log.Println("Planes:")
		for _, i := range keys {
			log.Printf("i: %2d p1: %3d p2: %3d p3: %3d \n", i, r.planes[i].p1, r.planes[i].p2, r.planes[i].p3)
		}
	}
	return nil
}

func (r *Rotation) Render(a *ui.Area, db *ui.AreaDrawParams) {
	log.Println("Rendering")
 	brush := ui.Brush{Type: ui.Solid, R: 0x00, G: 0x00, B: 0x00, A: 0x3F}
	strokeParams := ui.StrokeParams{
		Thickness: 1,
	}

	// Window centre
	wx := db.AreaWidth / 2
	wy := db.AreaHeight / 2

	path := ui.NewPath(ui.Winding)

	for i := range R.planes {

		path.NewFigure(
			R.points[R.planes[i].p1].x+wx,
			R.points[R.planes[i].p1].y+wy,
		)
		path.LineTo(
			R.points[R.planes[i].p2].x+wx,
			R.points[R.planes[i].p2].y+wy,
		)
		path.LineTo(
			R.points[R.planes[i].p3].x+wx,
			R.points[R.planes[i].p3].y+wy,
		)
		path.CloseFigure()

	}

	path.End()

	db.Context.Stroke(path, &brush, &strokeParams)
	db.Context.Save()
	db.Context.Restore()

	path.Free()
}
func (r *Rotation) Rotate() {

}

func main() {
	R = NewRotation()

	err = R.ReadOpts()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = R.ReadFile()
	if err != nil {
		log.Fatal(err.Error())
	}

	err := ui.Main(func() {
		window := ui.NewWindow("Rotation", 800, 600, false)
		window.OnClosing(func(*Window) bool {
			Quit()
			return true
		})
		WindowHandler := WindowHandler{Window: window}

		area := ui.NewArea(WindowHandler)
		area.Enable()
		area.Show()

		area.Handle()
		box := ui.NewVerticalBox()
		box.Append(area, true)
		//area.QueueRedrawAll()

		window.SetChild(box)
		window.Show()

 	})
	if err != nil {
		panic(err)
	}

}
