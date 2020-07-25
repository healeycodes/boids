package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/healeycodes/boids/vector"
)

type Vector2D = vector.Vector2D

const (
	screenWidth          = 1000
	screenHeight         = 1000
	numBoids             = 75
	maxForce             = 1.0
	maxSpeed             = 4.0
	alignPerception      = 75.0
	cohesionPerception   = 100.0
	separationPerception = 50.0
)

var (
	birdImage *ebiten.Image
)

func init() {
	fish, _, err := ebitenutil.NewImageFromFile("fish/chevron-up.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	w, h := fish.Size()
	birdImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 1)
	birdImage.DrawImage(fish, op)
}

type Boid struct {
	imageWidth  int
	imageHeight int
	pos         Vector2D
	vel         Vector2D
	acc         Vector2D
	angle       float64
}

func (boid *Boid) Rules(restOfFlock []*Boid) {
	alignSteering := Vector2D{}
	alignTotal := 0
	cohesionSteering := Vector2D{}
	cohesionTotal := 0
	separationSteering := Vector2D{}
	separationTotal := 0

	for i := range restOfFlock {
		other := restOfFlock[i]
		d := boid.pos.Distance(other.pos)
		if boid != other {
			if d < alignPerception {
				alignTotal++
				alignSteering.Add(other.vel)
			}
			if d < cohesionPerception {
				cohesionTotal++
				cohesionSteering.Add(other.pos)
			}
			if d < separationPerception {
				separationTotal++
				diff := boid.pos
				diff.Subtract(other.pos)
				diff.Divide(d)
				separationSteering.Add(diff)
			}
		}
	}

	if separationTotal > 0 {
		separationSteering.Divide(float64(separationTotal))
		separationSteering.SetMagnitude(maxSpeed)
		separationSteering.Subtract(boid.vel)
		separationSteering.SetMagnitude(maxForce * 1.2)
	}
	if cohesionTotal > 0 {
		cohesionSteering.Divide(float64(cohesionTotal))
		cohesionSteering.Subtract(boid.pos)
		cohesionSteering.SetMagnitude(maxSpeed)
		cohesionSteering.Subtract(boid.vel)
		cohesionSteering.SetMagnitude(maxForce)
	}
	if alignTotal > 0 {
		alignSteering.Divide(float64(alignTotal))
		alignSteering.SetMagnitude(maxSpeed)
		alignSteering.Subtract(boid.vel)
		alignSteering.Limit(maxForce)
	}

	boid.acc.Add(alignSteering)
	boid.acc.Add(cohesionSteering)
	boid.acc.Add(separationSteering)
}

func (boid *Boid) Update() {
	boid.angle = -1*math.Atan2(boid.vel.Y*-1, boid.vel.X) + math.Pi/2
	boid.pos.Add(boid.vel)
	boid.vel.Add(boid.acc)
	boid.vel.Limit(maxSpeed)
	boid.acc.Multiply(0.0)
}

func (boid *Boid) Edges() {
	if boid.pos.X < 0 {
		boid.pos.X = screenWidth
	} else if boid.pos.X > screenWidth {
		boid.pos.X = 0
	}
	if boid.pos.Y < 0 {
		boid.pos.Y = screenHeight
	} else if boid.pos.Y > screenHeight {
		boid.pos.Y = 0
	}
}

type Flock struct {
	boids []*Boid
}

func (flock *Flock) Update() {
	for i := range flock.boids {
		boid := flock.boids[i]
		boid.Edges()
		boid.Rules(flock.boids)
		boid.Update()
	}
}

type Game struct {
	flock  Flock
	op     ebiten.DrawImageOptions
	inited bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	rand.Seed(time.Hour.Milliseconds())
	g.flock.boids = make([]*Boid, numBoids)
	for i := range g.flock.boids {
		w, h := birdImage.Size()
		x, y := rand.Float64()*float64(screenWidth-w), rand.Float64()*float64(screenWidth-h)
		min, max := -maxForce, maxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		g.flock.boids[i] = &Boid{
			imageWidth:  w,
			imageHeight: h,
			pos:         Vector2D{X: x, Y: y},
			vel:         Vector2D{X: vx, Y: vy},
			acc:         Vector2D{X: 0, Y: 0},
		}
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !g.inited {
		g.init()
	}

	g.flock.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xff, 0xff, 0xff, 0xff})
	w, h := birdImage.Size()
	for i := range g.flock.boids {
		boid := g.flock.boids[i]
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		g.op.GeoM.Rotate(boid.angle)
		g.op.GeoM.Translate(boid.pos.X, boid.pos.Y)
		screen.DrawImage(birdImage, &g.op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
