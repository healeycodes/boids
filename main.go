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
	screenWidth  = 1024
	screenHeight = 768
	numBoids     = 300
	maxForce     = 1.0
	maxSpeed     = 4.0
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

func (v *Vector2D) Limit(z float64) {
	v.x = math.Max(-z, math.Min(v.x, z))
	v.y = math.Max(-z, math.Min(v.y, z))
}

func (v *Vector2D) Normalize() {
	mag := math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
	v.x /= mag
	v.y /= mag
}

func (v *Vector2D) SetMagnitude(z float64) {
	v.Normalize()
	v.x *= z
	v.y *= z
}

func (v *Vector2D) Divide(z float64) {
	v.x /= z
	v.y /= z
}

func (v *Vector2D) Multiply(z float64) {
	v.x *= z
	v.y *= z
}

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func (v Vector2D) Distance(v2 Vector2D) float64 {
	first := math.Pow(v2.x-v.x, 2)
	second := math.Pow(v2.y-v.y, 2)
	return math.Sqrt(first + second)
}

type Boid struct {
	imageWidth  int
	imageHeight int
	pos         Vector2D
	vel         Vector2D
	acc         Vector2D
}

func (boid *Boid) Align(restOfFlock []*Boid) {
	perception := 75.0
	steering := Vector2D{}
	total := 0
	for i := range restOfFlock {
		other := restOfFlock[i]
		d := boid.pos.Distance(other.pos)
		if boid != other && d < perception {
			total++
			steering.Add(other.vel)
		}
	}
	if total > 0 {
		steering.Divide(float64(total))
		steering.SetMagnitude(maxSpeed)
		steering.Subtract(boid.vel)
		steering.Limit(maxForce)
	}
	boid.acc.Add(steering)
}

func (boid *Boid) Cohesion(restOfFlock []*Boid) {
	perception := 100.0
	steering := Vector2D{}
	total := 0
	for i := range restOfFlock {
		other := restOfFlock[i]
		d := boid.pos.Distance(other.pos)
		if boid != other && d < perception {
			total++
			steering.Add(other.pos)
		}
	}
	if total > 0 {
		steering.Divide(float64(total))
		steering.Subtract(boid.pos)
		steering.SetMagnitude(maxSpeed)
		steering.Subtract(boid.vel)
		steering.SetMagnitude(maxForce)
	}
	boid.acc.Add(steering)
}

func (boid *Boid) Separation(restOfFlock []*Boid) {
	perception := 50.0
	steering := Vector2D{}
	total := 0
	for i := range restOfFlock {
		other := restOfFlock[i]
		d := boid.pos.Distance(other.pos)
		if boid != other && d < perception {
			total++
			diff := boid.pos
			diff.Subtract(other.pos)
			diff.Divide(d)
			steering.Add(diff)
		}
	}
	if total > 0 {
		steering.Divide(float64(total))
		steering.SetMagnitude(maxSpeed)
		steering.Subtract(boid.vel)
		steering.SetMagnitude(maxForce * 1.1)
	}
	boid.acc.Add(steering)
}

func (boid *Boid) Update() {
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
	// snapspot := make([]*Boid, len(flock.sprites))
	// for i := range flock.sprites {
	// 	boid := flock.sprites[i]
	// 	snapspot[i] = &Boid{pos: boid.pos, vel: boid.vel, acc: boid.acc}
	// }
	for i := range flock.sprites {
		boid := flock.sprites[i]
		boid.Edges()
		boid.Align(flock.sprites)
		boid.Cohesion(flock.sprites)
		boid.Separation(flock.sprites)
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
		// g.op.GeoM.Rotate(boid.angle * math.Pi / 180)
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
