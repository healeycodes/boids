# Boids with Go and Ebiten

_Under construction_ 👷🏻‍♀️

👉🏻  [Demo link](https://healeycodes.github.io/boids/) to the WASM version.

![Animated GIF of a flocking simulation](https://github.com/healeycodes/boids/raw/master/dist/preview.gif)

### TODOs

- Zero optimizations
- A little messy
- Jittery (don't know why yet)
- Doesn't use a snapshot for each 'generation'
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
