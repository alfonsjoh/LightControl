build:
	go build -ldflags -H=windowsgui -o ./out/LightControl.exe ./src/

run:
	go run ./src/
