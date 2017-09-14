package main

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andlabs/ui"
	"github.com/davecgh/go-spew/spew"
	flag "github.com/ogier/pflag"
)

var lightSource = point{
	0,
	0,
	100,
}

var err error
var R *Rotation
var mutex = &sync.Mutex{}

// A point is a point in three dimensional space
type point struct {
	x float64
	y float64
	z float64
}

// A plane is an ordered set of points, each of which in an index to the point collection. The points are ordered such
// that p1->p2 forms vector A and points p2->p3 form vector B. With the points being ordered it is possible to calculate
// a vector C that is the orthogonal normal to the plane.
type plane struct {
	p1    int
	p2    int     // p1->p2 = vA
	p3    int     // p2->p3 = vB
	theta float64 // the angle between the orthogonal normal of the plane defined by vA->vB and the vector from the origin to the light source.
}

type Config struct {
	debug     int
	filename  string
	wireframe bool
	zoom      float64
	delay     int64
	// Rotation factors
	xy, xz, yz float64
}

type Rotation struct {
	pointCount int
	planeCount int
	points     map[int]point
	planes     map[int]plane
	config     Config
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
			1.0,
			100,
			.03,
			.06,
			.09,
		},
	}
	return &o
}

type Rotatable interface {
	ReadOpts() error
	ReadFile() error
	Rotate()
}

func (r *Rotation) ReadOpts() error {
	flag.StringVarP(&r.config.filename, "filename", "f", "", "Data file to load")
	flag.Int64VarP(&r.config.delay, "delay", "d", 100, "Delay (ms) to add to rotation. (Default: 100ms) ")
	flag.BoolVarP(&r.config.wireframe, "wireframe", "w", false, "Show object as a wirefram? (Default: no)")
	flag.IntVarP(&r.config.debug, "debug", "g", 0, "Debug logging level. 0=none, 1=some (Default: 0)")
	flag.Float64VarP(&r.config.zoom, "zoom", "z", 1.0, "Zoom factor")
	flag.Float64Var(&r.config.xy, "xy", 0.1, "XY rotation factor")
	flag.Float64Var(&r.config.xz, "xz", 0.1, "XZ rotation factor")
	flag.Float64Var(&r.config.yz, "yz", 0.1, "YZ rotation factor")
	flag.Parse()

	if r.config.filename == "" {
		return errors.New("Missing filename")
	}

	if r.config.delay < 0 {
		return errors.New("Delay must be positive.")
	}

	if r.config.zoom < 0 {
		return errors.New("Zoom factor must be positive.")
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
			0.0,
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
			log.Printf("i: %2d x: %3f y: %3f z: %3f \n", i, r.points[i].x, r.points[i].y, r.points[i].z)
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

func (r *Rotation) Rotate() {
	for i := range r.points {
		var x0, y0, z0 float64
		var x1, y1, z1 float64
		var x2, y2, z2 float64
		var x3, y3, z3 float64

		x0 = r.points[i].x
		y0 = r.points[i].y
		z0 = r.points[i].z

		x1 = x0*math.Cos(r.config.xy) - y0*math.Sin(r.config.xy)
		y1 = x0*math.Sin(r.config.xy) + y0*math.Cos(r.config.xy)
		z1 = z0

		x2 = x1*math.Cos(r.config.xz) - z1*math.Sin(r.config.xz)
		y2 = y1
		z2 = x1*math.Sin(r.config.xz) + z1*math.Cos(r.config.xz)

		x3 = x2
		y3 = y2*math.Cos(r.config.yz) - z2*math.Sin(r.config.yz)
		z3 = y2*math.Sin(r.config.yz) + z2*math.Cos(r.config.yz)

		mutex.Lock()
		r.points[i] = point{x3, y3, z3}
		mutex.Unlock()
	}

	for i := range r.planes {
		var a1, b1, c1 float64
		var a2, b2, c2 float64
		var a3, b3, c3 float64
		var theta float64

		mutex.Lock()
		a1 = r.points[r.planes[i].p2].x - r.points[r.planes[i].p1].x
		b1 = r.points[r.planes[i].p2].y - r.points[r.planes[i].p1].y
		c1 = r.points[r.planes[i].p2].z - r.points[r.planes[i].p1].z

		a2 = r.points[r.planes[i].p3].x - r.points[r.planes[i].p2].x
		b2 = r.points[r.planes[i].p3].y - r.points[r.planes[i].p2].y
		c2 = r.points[r.planes[i].p3].z - r.points[r.planes[i].p2].z
		mutex.Unlock()

		a3 = b1*c2 - c1*b2
		b3 = c1*a2 - a1*c2
		c3 = a1*b2 - b1*a2

		dotProduct := ((a3 * lightSource.x) + (b3 * lightSource.y) + (c3 * lightSource.z))
		lensq1 := (a3*a3 + b3*b3 + c3*c3)
		lensq2 := (lightSource.x*lightSource.x + lightSource.y*lightSource.y + lightSource.z*lightSource.z)
		theta = dotProduct / math.Sqrt(lensq1*lensq2)

		if r.config.debug > 2 {
			log.Printf("Plane %d: %f\n", i, theta)
		}

		mutex.Lock()
		p := r.planes[i]
		p.theta = theta
		r.planes[i] = p
		mutex.Unlock()

	}
}

func (r *Rotation) Render(a *ui.Area, db *ui.AreaDrawParams) {
	strokeParams := ui.StrokeParams{
		Thickness: 1,
	}

	// Window centre
	wx := db.AreaWidth / 2
	wy := db.AreaHeight / 2

	mutex.Lock()
	for i := range r.planes {
		if r.planes[i].theta > 0 || r.config.wireframe {
			path := ui.NewPath(ui.Winding)
			path.NewFigure(
				r.points[r.planes[i].p1].x*r.config.zoom+wx,
				r.points[r.planes[i].p1].y*r.config.zoom+wy,
			)
			path.LineTo(
				r.points[r.planes[i].p2].x*r.config.zoom+wx,
				r.points[r.planes[i].p2].y*r.config.zoom+wy,
			)
			path.LineTo(
				r.points[r.planes[i].p3].x*r.config.zoom+wx,
				r.points[r.planes[i].p3].y*r.config.zoom+wy,
			)
			path.CloseFigure()
			path.End()

			if r.config.wireframe {
				db.Context.Stroke(path, &ui.Brush{Type: ui.Solid, R: 0x00, G: 0x00, B: 0x00, A: 0x3F}, &strokeParams)
			} else {
				db.Context.Fill(path, &ui.Brush{Type: ui.Solid, R: r.planes[i].theta, G: r.planes[i].theta, B: r.planes[i].theta, A: 0xFF})
			}
		}
	}
	mutex.Unlock()

	db.Context.Save()
	db.Context.Restore()
}

func (r *Rotation) Run(area *ui.Area) {
	if r.config.debug > 2 {
		spew.Dump("delay: ", r.config.delay)
	}
	for {
		r.Rotate()

		time.Sleep(time.Duration(r.config.delay) * time.Millisecond)
		ui.QueueMain(func() {
			area.QueueRedrawAll()
		})
	}
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
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
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

		go R.Run(area)

	})
	if err != nil {
		panic(err)
	}

}
