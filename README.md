## Boids with Go and Ebiten

> My blog post: [Boids in WebAssembly Using Go](https://healeycodes.com/boids-flocking-simulation/)

<br>

👉🏻  [Demo link](https://healeycodes.github.io/boids/) to the WASM version.

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
./build.sh
```

This compiles the program into WebAssembly and gzips the result.

Why gzip? We can shrink the file from 8MB to 2MB and include a 44KB library to inflate the file in a browser.

<br>

License: MIT.
