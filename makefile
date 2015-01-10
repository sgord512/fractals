all: fractals

fractals: fractals.go poly.go glutil.go
	go build -o $@ $^ 
