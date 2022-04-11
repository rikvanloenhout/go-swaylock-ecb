# go-swaylock-ecb
AES-128 ECB mode encrypted screenshots for swaylock. 

`go-swaylock-ecb` takes screenshots of your monitor output and encrypts them with AES-128 in ECB mode (known for the ECB penguin) and open `swaylock` with the encrypted output as a background.

For more information on how this works and why it is cool, please read [The ECB Penguin by Filippo Valsorda](https://words.filippo.io/the-ecb-penguin/).

> :warning: Please note that this program's function is purely aesthetical. Encrypting the background of your lockscreen does not increase privacy or security.

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
