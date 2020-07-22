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
)

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
	fishImage *ebiten.Image
)

func init() {
	fish, _, err := ebitenutil.NewImageFromFile("fish/chevron-up.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	w, h := fish.Size()
	fishImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 1)
	fishImage.DrawImage(fish, op)
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
	boid.angle = -1*math.Atan2(boid.vel.y*-1, boid.vel.x) + math.Pi/2
	boid.pos.Add(boid.vel)
	boid.vel.Add(boid.acc)
	boid.vel.Limit(maxSpeed)
	boid.acc.Multiply(0.0)
}

func (boid *Boid) Edges() {
	if boid.pos.x < 0 {
		boid.pos.x = screenWidth
	} else if boid.pos.x > screenWidth {
		boid.pos.x = 0
	}
	if boid.pos.y < 0 {
		boid.pos.y = screenHeight
	} else if boid.pos.y > screenHeight {
		boid.pos.y = 0
	}
}

type Boids struct {
	sprites []*Boid
}

func (flock *Boids) Update() {
	for i := range flock.sprites {
		boid := flock.sprites[i]
		boid.Edges()
		boid.Rules(flock.sprites)
		boid.Update()
	}
}

type Game struct {
	flock  Boids
	op     ebiten.DrawImageOptions
	inited bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	rand.Seed(time.Hour.Milliseconds())
	g.flock.sprites = make([]*Boid, numBoids)
	for i := range g.flock.sprites {
		w, h := fishImage.Size()
		x, y := rand.Float64()*float64(screenWidth-w), rand.Float64()*float64(screenWidth-h)
		min, max := -2.0, 2.0
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		g.flock.sprites[i] = &Boid{
			imageWidth:  w,
			imageHeight: h,
			pos:         Vector2D{x: x, y: y},
			vel:         Vector2D{x: vx, y: vy},
			acc:         Vector2D{x: 0, y: 0},
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
	w, h := fishImage.Size()
	for i := range g.flock.sprites {
		boid := g.flock.sprites[i]
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		g.op.GeoM.Rotate(boid.angle)
		g.op.GeoM.Translate(boid.pos.x, boid.pos.y)
		screen.DrawImage(fishImage, &g.op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// TODO: split into a vector library

type Vector2D struct {
	x float64
	y float64
}

func (v *Vector2D) Add(v2 Vector2D) {
	v.x += v2.x
	v.y += v2.y
}

func (v *Vector2D) Subtract(v2 Vector2D) {
	v.x -= v2.x
	v.y -= v2.y
}

func (v *Vector2D) Limit(max float64) {
	magSq := v.MagnitudeSquared()
	if magSq > max*max {
		v.Divide(math.Sqrt(magSq))
		v.Multiply(max)
	}
}

func (v *Vector2D) Normalize() {
	mag := math.Sqrt(v.x*v.x + v.y*v.y)
	v.x /= mag
	v.y /= mag
}

func (v *Vector2D) SetMagnitude(z float64) {
	v.Normalize()
	v.x *= z
	v.y *= z
}

func (v *Vector2D) MagnitudeSquared() float64 {
	return v.x*v.x + v.y*v.y
}

func (v *Vector2D) Divide(z float64) {
	v.x /= z
	v.y /= z
}

func (v *Vector2D) Multiply(z float64) {
	v.x *= z
	v.y *= z
}

func (v Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v2.x-v.x, 2) + math.Pow(v2.y-v.y, 2))
}
