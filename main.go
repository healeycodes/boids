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
	imageWidth   int
	imageHeight  int
	position     Vector2D
	velocity     Vector2D
	acceleration Vector2D
}

func (boid *Boid) ApplyRules(restOfFlock []*Boid) {
	alignSteering := Vector2D{}
	alignTotal := 0
	cohesionSteering := Vector2D{}
	cohesionTotal := 0
	separationSteering := Vector2D{}
	separationTotal := 0

	for _, other := range restOfFlock {
		d := boid.position.Distance(other.position)
		if boid != other {
			if d < alignPerception {
				alignTotal++
				alignSteering.Add(other.velocity)
			}
			if d < cohesionPerception {
				cohesionTotal++
				cohesionSteering.Add(other.position)
			}
			if d < separationPerception {
				separationTotal++
				diff := boid.position
				diff.Subtract(other.position)
				diff.Divide(d)
				separationSteering.Add(diff)
			}
		}
	}

	if separationTotal > 0 {
		separationSteering.Divide(float64(separationTotal))
		separationSteering.SetMagnitude(maxSpeed)
		separationSteering.Subtract(boid.velocity)
		separationSteering.SetMagnitude(maxForce * 1.2)
	}
	if cohesionTotal > 0 {
		cohesionSteering.Divide(float64(cohesionTotal))
		cohesionSteering.Subtract(boid.position)
		cohesionSteering.SetMagnitude(maxSpeed)
		cohesionSteering.Subtract(boid.velocity)
		cohesionSteering.SetMagnitude(maxForce * 0.9)
	}
	if alignTotal > 0 {
		alignSteering.Divide(float64(alignTotal))
		alignSteering.SetMagnitude(maxSpeed)
		alignSteering.Subtract(boid.velocity)
		alignSteering.Limit(maxForce)
	}

	boid.acceleration.Add(alignSteering)
	boid.acceleration.Add(cohesionSteering)
	boid.acceleration.Add(separationSteering)
	boid.acceleration.Divide(3)
}

func (boid *Boid) ApplyMovement() {
	boid.position.Add(boid.velocity)
	boid.velocity.Add(boid.acceleration)
	boid.velocity.Limit(maxSpeed)
	boid.acceleration.Multiply(0.0)
}

func (boid *Boid) CheckEdges() {
	if boid.position.X < 0 {
		boid.position.X = screenWidth
	} else if boid.position.X > screenWidth {
		boid.position.X = 0
	}
	if boid.position.Y < 0 {
		boid.position.Y = screenHeight
	} else if boid.position.Y > screenHeight {
		boid.position.Y = 0
	}
}

type Flock struct {
	boids []*Boid
}

func (flock *Flock) Logic() {
	for _, boid := range flock.boids {
		boid.CheckEdges()
		boid.ApplyRules(flock.boids)
		boid.ApplyMovement()
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
			imageWidth:   w,
			imageHeight:  h,
			position:     Vector2D{X: x, Y: y},
			velocity:     Vector2D{X: vx, Y: vy},
			acceleration: Vector2D{X: 0, Y: 0},
		}
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !g.inited {
		g.init()
	}

	g.flock.Logic()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	w, h := birdImage.Size()
	for _, boid := range g.flock.boids {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		g.op.GeoM.Rotate(-1*math.Atan2(boid.velocity.Y*-1, boid.velocity.X) + math.Pi/2)
		g.op.GeoM.Translate(boid.position.X, boid.position.Y)
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
