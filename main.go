package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"image"
	"math"
	"os/exec"

	"github.com/Difrex/gosway/ipc"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"github.com/fogleman/gg"
	"github.com/lmittmann/ppm"
)

func main() {
	// Generate 16-bit key for AES-128
	key := make([]byte, 16)
	rand.Read(key)

	// Get sway outputs
	outputs, err := outputs()
	if err != nil {
		panic(err)
	}

	screens := map[string]string{}

	for _, output := range outputs {

		// Get screenshot in PPM format
		pt, err := grim(output)
		if err != nil {
			panic(err)
		}

		// Encrypt the file data (and skip the file header)
		ct, err := encrypt(pt, key)
		if err != nil {
			panic(err)
		}
		// ct := pt

		// Create an image.Image from the PPM file
		img, err := ppm.Decode(bytes.NewReader(ct))
		if err != nil {
			panic(err)
		}

		context := gg.NewContextForImage(img)

		// Apply pixelate effect
		pixelate(img, context, 100)

		// Path for file in tmp dir
		path := "/tmp/" + output.Name + ".png"

		// Write file to tmp dir
		err = context.SavePNG(path)
		if err != nil {
			panic(err)
		}

		screens[output.Name] = path
	}

	// Exec swaylock
	if err = swaylock(screens); err != nil {
		panic(err)
	}
}

func encrypt(src, key []byte) ([]byte, error) {
	divider := []byte("\n")
	parts := bytes.SplitN(src, divider, 5)
	pt := parts[4]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt)
	if err != nil {
		return nil, err
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	parts[4] = ct
	return bytes.Join(parts, divider), nil
}

func outputs() ([]*ipc.Output, error) {
	conn, err := ipc.NewSwayConnection()
	if err != nil {
		return nil, err
	}

	outputs, err := conn.GetOutputs()
	if err != nil {
		return nil, err
	}

	activeOutputs := make([]*ipc.Output, 0)
	for _, output := range outputs {
		if output.Active {
			activeOutputs = append(activeOutputs, output)
		}
	}

	return activeOutputs, nil
}

func grim(output *ipc.Output) ([]byte, error) {
	return exec.Command("grim", "-t", "ppm", "-o", output.Name, "-").Output()
}

func swaylock(outputs map[string]string) error {
	params := []string{"-f", "-e", "--indicator-radius", "85"}

	for output, file := range outputs {
		params = append(params, "-i", output+":"+file)
	}

	return exec.Command("swaylock", params...).Run()
}

func pixelate(img image.Image, context *gg.Context, factor float64) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	wScale := int(math.Round(float64(width) / factor))
	hScale := int(math.Round(float64(height) / factor))
	w, h := 1, 1

	for y := 0; y < height; y += h {
		h = hScale
		for x := 0; x < width; x += w {
			w = wScale
			context.Push()
			if x+w > width {
				w = width - x
			}
			if y+h > height {
				h = height - y
			}
			context.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
			r, g, b := getAverageColor(x, y, x+w, y+h, img)

			context.SetRGB255(r, g, b)
			context.Fill()
			context.Pop()
		}
	}
}

// Returns the average red, green, and blue colors within a rectangle of the image
func getAverageColor(x0, y0, x1, y1 int, img image.Image) (int, int, int) {
	r := make([]uint32, (y1-y0)*(x1-x0))
	g := make([]uint32, (y1-y0)*(x1-x0))
	b := make([]uint32, (y1-y0)*(x1-x0))
	idx := 0

	// Get all colors in the range and add them to their slices
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			color := img.At(x, y)
			red, green, blue, _ := color.RGBA()
			r[idx], g[idx], b[idx] = red, green, blue
			idx += 1
		}
	}

	// Get the averages for each color slice
	dR := getAverage(r)
	dG := getAverage(g)
	dB := getAverage(b)
	return dR, dG, dB
}

// Takes a slice of uint32 and returns an average color in an int
func getAverage(slice []uint32) int {
	var sum uint32 = 0
	for i := 0; i < len(slice); i++ {
		sum += (slice[i])
	}
	count := len(slice)
	// Take the square root to get the true average
	avg := math.Sqrt(float64(sum / uint32(count)))

	return int(avg)
}
