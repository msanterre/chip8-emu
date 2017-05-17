# chip8-emulator

**Note: Work in progress**

## Purpose
I wanted to learn how emulators/chips worked, so I decided to make this small chip-8 emulator. 

## Method
I started by working off [Wikipedia](https://en.wikipedia.org/wiki/CHIP-8), but it was missing some components, 
so I ended up using [Cowgod's technical reference guide](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM).

## Graphics & keyboard
I'm using SDL to write to the screen and get keyboard input. It should work for all platforms, but if it doesn't, 
please create an issue and I'll investigate.

## Running
The normal way of running is:
```
go run main.go <rom_path>
````

The roms are in the directory, so you can just do: 
```
go run main.go roms/PONG
```

**If you're running go v1.8 and it's giving you `signal: killed`. Run this instead:**
```
go run -ldflags=-s main.go roms/PONG
```


## TODO
- [ ] Keyboard input seems to register, but it needs to be mapped properly on normal keyboards
- [ ] Fix out index out of range issue during the draw
- [ ] Rom guides
