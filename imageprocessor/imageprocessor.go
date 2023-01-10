package imageprocessor

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"sort"
)

func WriteResultsToCSV(url string, colors []color.Color, errorChan chan error) {
	// Open the CSV file
	file, err := os.OpenFile("output.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errorChan <- fmt.Errorf("can't open output.csv: %v", err)
		fmt.Println("can't open output.csv: ", err)
		return
	}
	defer file.Close()

	// Write the URL and dominant colors to the CSV file
	_, err = file.WriteString(url)
	for _, c := range colors {
		_, err = file.WriteString(",")
		r, g, b, _ := c.RGBA()
		_, err = file.WriteString(fmt.Sprintf("#%02x%02x%02x", r, g, b))
		fmt.Println("DONE! ", url, r, g, b)
	}
	_, err = file.WriteString("\n")
	if err != nil {
		errorChan <- fmt.Errorf("problem writing out to the file: %v", err)
		fmt.Println("problem writing out to the file: ", err)
		return
	}
	return
}

//Simplest approach to extracting the dominant colors.
//Counts the number of pixels of each color in the image, and
//return the 3 most common colors.
func ExtractDominantColors(img image.Image) []color.Color {

	//create a map for each color's count
	colorCount := make(map[color.Color]int)
	//loop through each pixel in the image and increment the count for the appropriate color
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colorCount[img.At(x, y)]++
		}
	}

	//sort our color counts
	var colors []color.Color
	for c := range colorCount {
		colors = append(colors, c)
	}

	//I could handcraft a more efficient sorting algorithm here, however
	//the compiler under the hood is likely pattern matching this code and
	//optimizing it better than I could. I'll leave it as is until I need
	//to change it for a specific goal.
	sort.Slice(colors, func(i, j int) bool {
		return colorCount[colors[i]] > colorCount[colors[j]]
	})

	// Return the 3 most prevalent colors
	fmt.Println("Dominant Colors: ", colors[0], colors[1], colors[2])
	return []color.Color{colors[0], colors[1], colors[2]}
}

//goland:noinspection SpellCheckingInspection
func extractDominantColorsKmeans(img image.Image) []color.Color {
	//initialize the centroids with random pixels
	pixels := getPixels(img)
	centroids := []color.Color{pixels[rand.Intn(len(pixels))], pixels[rand.Intn(len(pixels))], pixels[rand.Intn(len(pixels))]}

	//iterate until the centroids converge
	for {
		//assign each pixel to its nearest centroid
		clusters := [][]color.Color{{}, {}, {}}

		for _, p := range pixels {
			minimumDistance := math.MaxFloat64
			minimumIndex := 0
			for i, c := range centroids {
				distance := distance(p, c)
				if distance < minimumDistance {
					minimumDistance = distance
					minimumIndex = i
				}
			}
			clusters[minimumIndex] = append(clusters[minimumIndex], p)
		}
		// calculate the new centroids
		newCentroids := []color.Color{averageColor(clusters[0]), averageColor(clusters[1]), averageColor(clusters[2])}
		//have the centroids converged yet?
		if centroidsEqual(centroids, newCentroids) {
			fmt.Println("Centroids Converged!")
			return centroids
		}
		//if not, update the centroids
		centroids = newCentroids
	}
}

//return a slice of color.RGBA values for the pixels in a given image
func getPixels(img image.Image) []color.Color {
	//get the image dimensions
	bounds := img.Bounds()
	//initialize the pixel slice
	pixels := make([]color.Color, 0, bounds.Max.X*bounds.Max.Y)
	//iterate over the pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels = append(pixels, img.At(x, y))
		}
	}
	return pixels
}

//get the euclidean distance between two colors
func distance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	rInt := int32(r1) - int32(r2)
	gInt := int32(g1) - int32(g2)
	bInt := int32(b1) - int32(b2)

	return math.Sqrt(float64(rInt*rInt + gInt*gInt + bInt*bInt))
}

//return the average of a slice of colors.
func averageColor(colors []color.Color) color.Color {
	var r, g, b uint32
	for _, c := range colors {
		r1, g1, b1, _ := c.RGBA()
		r += r1
		g += g1
		b += b1
	}
	n := uint32(len(colors))
	if n > 0 {
		return color.RGBA{uint8(r / n), uint8(g / n), uint8(b / n), 255}

	} else {
		return color.RGBA{0, 0, 0, 0}
	}
}

//are the two centroids slices the same?
func centroidsEqual(c1, c2 []color.Color) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i, c := range c1 {
		if c != c2[i] {
			return false
		}
	}
	return true
}
