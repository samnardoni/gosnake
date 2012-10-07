package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"github.com/jteeuwen/glfw"
	"os"
)

type Counter struct {
	hit, current int
}

func (c *Counter) Tick() bool {
	c.current++
	if c.current%c.hit == 0 {
		return true
	}
	return false
}

type Coord struct {
	x, y int
}

const (
	Title  = "Snake"
	Width  = 480
	Height = 480
)

const (
	DirUp = iota
	DirDown
	DirLeft
	DirRight
)

var (
	snake   []Coord
	food    map[Coord]bool = make(map[Coord]bool)
	grow    map[Coord]bool = make(map[Coord]bool)
	dir     int            = DirUp
	running bool           = true
	counter Counter        = Counter{hit: 6, current: 0}
)

func main() {
	var err error

	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.Terminate()

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	initGL()
	initSnake()
	initFood()

	for running && glfw.WindowParam(glfw.Opened) == 1 {
		update()
		drawScene()
	}
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onKey(key, state int) {
	// TODO Magic 1?
	if state != 1 {
		return
	}

	switch key {
	case glfw.KeyEsc:
		running = false
	case glfw.KeyUp:
		dir = DirUp
	case glfw.KeyDown:
		dir = DirDown
	case glfw.KeyLeft:
		dir = DirLeft
	case glfw.KeyRight:
		dir = DirRight
	}
}

func initGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

func initSnake() {
	snake = append(snake, Coord{x: 1, y: 2})
	snake = append(snake, Coord{x: 2, y: 2})
	snake = append(snake, Coord{x: 2, y: 3})
	snake = append(snake, Coord{x: 2, y: 4})
	snake = append(snake, Coord{x: 3, y: 4})
	snake = append(snake, Coord{x: 4, y: 4})
	dir = DirUp
}

func initFood() {
	food[Coord{x: 3, y: 3}] = true
	food[Coord{x: 7, y: 3}] = true
	food[Coord{x: 3, y: 5}] = true
	food[Coord{x: 6, y: 3}] = true
}

func update() {

	if hit := counter.Tick(); !hit {
		return
	}

	head := snake[len(snake)-1]
	tail := snake[0]

	// Eat food
	if _, present := food[head]; present {
		delete(food, head)
		grow[tail] = true
	}

	// Move snake
	switch dir {
	case DirUp:
		snake = append(snake, Coord{x: head.x, y: head.y + 1})
	case DirDown:
		snake = append(snake, Coord{x: head.x, y: head.y - 1})
	case DirLeft:
		snake = append(snake, Coord{x: head.x - 1, y: head.y})
	case DirRight:
		snake = append(snake, Coord{x: head.x + 1, y: head.y})
	}

	if _, present := grow[tail]; present {
		delete(grow, tail)
	} else {
		snake = snake[1:]
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Food
	for coord, _ := range food {
		gl.LoadIdentity()
		gl.Scalef(0.05, 0.05, 1)
		gl.Translatef(float32(coord.x)-5, float32(coord.y)-5, -1)

		gl.Begin(gl.QUADS)
		gl.Color3f(0.5, float32(coord.x)/10, float32(coord.y)/10)
		gl.Vertex3f(0.25, 0.75, 0)
		gl.Vertex3f(0.75, 0.75, 0)
		gl.Color3f(0.3, float32(coord.x)/10-0.2, float32(coord.y)/10-0.2)
		gl.Vertex3f(0.75, 0.25, 0)
		gl.Vertex3f(0.25, 0.25, 0)
		gl.End()
	}

	//Snake
	for _, coord := range snake {
		gl.LoadIdentity()
		gl.Scalef(0.05, 0.05, 1)
		gl.Translatef(float32(coord.x)-5, float32(coord.y)-5, -1)

		gl.Begin(gl.QUADS)
		gl.Color3f(float32(coord.y)/10, float32(coord.x)/10, 0.5)
		gl.Vertex3f(0, 1, 0)
		gl.Vertex3f(1, 1, 0)
		gl.Color3f(float32(coord.y)/10-0.2, float32(coord.x)/10-0.2, 0.3)
		gl.Vertex3f(1, 0, 0)
		gl.Vertex3f(0, 0, 0)
		gl.End()
	}

	glfw.SwapBuffers()
}
