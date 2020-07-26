GOOS=js GOARCH=wasm go build -o dist/boids.wasm github.com/healeycodes/boids
gzip -9 -v -c dist/boids.wasm > dist/boids.wasm.gz
