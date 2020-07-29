## Boids with Go and Ebiten

> My blog post: [Boids in WebAssembly Using Go](https://healeycodes.com/boids-flocking-simulation/)

<br>

👉🏻 [Demo link](https://healeycodes.github.io/boids/) to the WASM version.

I wrote this program, an implementation of Craig Reynolds' _Boids_, in order to learn more about deploying Go on the web, and to tackle a problem that escaped me when I was learning to code!

<br>

![Animated GIF of a flocking simulation](https://github.com/healeycodes/boids/raw/master/dist/preview.gif)

### Possible improvements

- Field of vision support (boids shouldn't look behind 👀)
- QuadTree optimization
- Different `maxSpeed`/`maxForce` for each boid
- Graphical interface for live-editing of values
- Use a snapshot for each 'generation'
- ..

<br>

### Run 🦢

```
go run main.go
```

### Build 🕊

```
GOOS=js GOARCH=wasm go build -o dist/boids.wasm github.com/healeycodes/boids
```

This compiles the program into WebAssembly in the `dist` folder.

The simulation can be viewed in a web browser via `index.html`. To get this working locally, you may need to serve the files from a web server.

<br>

License: MIT.
