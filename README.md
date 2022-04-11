# go-swaylock-ecb
AES-128 ECB mode encrypted screenshots for all outputs as a background for swaylock. 

When executed, the program will take screenshots using `grim`, encrypt them AES-128 in ECB mode and run `swaylock` with the encrypted screenshots as wallpapers for each output.

For more information on how this works and why it is cool, please read [The ECB Penguin by Filippo Valsorda](https://words.filippo.io/the-ecb-penguin/).

## Dependencies
- [go](https://github.com/golang/go)
- [swaylock](https://github.com/swaywm/swaylock)
- [grim](https://sr.ht/~emersion/grim/)

## Installation
Run `go install github.com/rikvanloenhout/go-swaylock-ecb@latest`.

## Usage
Run `go-swaylock-ecb` to lock your screen. 

## TODOs
- Extinguish vertical line pattern
- Improve performance (drastically)

## Example sceenshot
![AES-ECB encrypted screenshot of a 4k output](examples/screenshot.png "4k AES-EDB encrypted screenshot")
