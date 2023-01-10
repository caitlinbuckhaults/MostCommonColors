package imageprocessor

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"reflect"
	"testing"
)

func TestProcessImage(t *testing.T) {
	// TODO: DO THIS
}

func TestExtractDominantColors(t *testing.T) {
	//TODO: DO THIS
}

func TestDownloadAndProcessImage(t *testing.T) {
	type args struct {
		url        string
		resultChan chan []color.Color
		errorChan  chan error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DownloadAndProcessImage(tt.args.url, tt.args.resultChan, tt.args.errorChan)
		})
	}
}

func TestWriteResultsToCSV(t *testing.T) {
	type args struct {
		url       string
		colors    []color.Color
		errorChan chan error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteResultsToCSV(tt.args.url, tt.args.colors, tt.args.errorChan)
		})
	}
}

func Test_averageColor(t *testing.T) {
	type args struct {
		colors []color.Color
	}
	tests := []struct {
		name string
		args args
		want color.Color
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := averageColor(tt.args.colors); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("averageColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_centroidsEqual(t *testing.T) {
	type args struct {
		c1 []color.Color
		c2 []color.Color
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := centroidsEqual(tt.args.c1, tt.args.c2); got != tt.want {
				t.Errorf("centroidsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeJpegFiles(t *testing.T) {

	testCases := []struct {
		name     string
		fileName string
		wantErr  error
		wantType image.Image
	}{
		{
			name:     "Valid image",
			fileName: "validImage.jpeg",
			wantErr:  nil,
			wantType: &image.RGBA{},
		},
		{
			name:     "Invalid image format",
			fileName: "invalidFormat.webp",
			wantErr:  fmt.Errorf("unexpected EOI"),
			wantType: nil,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeJpeg(tc.fileName)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error: got %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && err.Error() != tc.wantErrStr {
				t.Errorf("unexpected error string: got %q, want %q", err.Error(), tc.wantErrStr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDecodeJpegBytes(t *testing.T) {
	// Define the test data
	validImageData := []byte{}
	invalidImageFormatData := []byte{
		0x47, 0x49, 0x46, 0x38, 0x37, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00, 0x00, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
		0x01, 0x00, 0x3B,
	}
	invalidImageData := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x01, 0x00, 0x48,
		0x00, 0x48, 0x00, 0x00, 0xFF, 0xE1, 0x00, 0x68, 0x45, 0x78, 0x69, 0x66, 0x00, 0x00, 0x49, 0x49,
		0x2A, 0x00, 0x08, 0x00, 0x00, 0x00, 0x04, 0x00, 0x1A, 0x01, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x3E, 0x00, 0x00, 0x00, 0x1B, 0x01, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00, 0x46, 0x00, 0x00, 0x00,
		0x28, 0x01, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x31, 0x01, 0x02, 0x00}
	testCases := []struct {
		name       string
		reader     io.Reader
		want       image.Image
		wantErr    bool
		wantErrStr string
	}{
		{
			name:       "valid image",
			reader:     bytes.NewReader(validImageData),
			want:       validImage,
			wantErr:    false,
			wantErrStr: "",
		},
		{
			name:       "invalid image format",
			reader:     bytes.NewReader(invalidImageFormatData),
			want:       nil,
			wantErr:    true,
			wantErrStr: "image: unknown format",
		},
		{
			name:       "invalid image data",
			reader:     bytes.NewReader(invalidImageData),
			want:       nil,
			wantErr:    true,
			wantErrStr: "image: invalid format: invalid JPEG format: missing SOI marker",
		},
		// Add more test cases here
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeJpeg(tc.reader)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error: got %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && err.Error() != tc.wantErrStr {
				t.Errorf("unexpected error string: got %q, want %q", err.Error(), tc.wantErrStr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func Test_distance(t *testing.T) {
	tests := []struct {
		name string
		c1   color.Color
		c2   color.Color
		want float64
	}{
		{"Distance RGB Zero", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 255}, 0},
		{"Distance RGB 255", color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, 0},
		{"Distance RGB Max", color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}, 441.67295593},
		{"Distance RG Value Swap", color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, 358.309518948},
		{"Distance GB Value Swap", color.RGBA{0, 0, 255, 255}, color.RGBA{0, 255, 0, 255}, 170.7038838},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := distance(tt.c1, tt.c2)
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("distance(%v, %v) = %v, want %v", tt.c1, tt.c2, got, tt.want)
			}
		})
	}
}

func Test_extractDominantColors(t *testing.T) {
	type args struct {
		img image.Image
	}
	tests := []struct {
		name string
		args args
		want []color.Color
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractDominantColors(tt.args.img); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractDominantColors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractDominantColorsKmeans(t *testing.T) {
	type args struct {
		img image.Image
	}
	tests := []struct {
		name string
		args args
		want []color.Color
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractDominantColorsKmeans(tt.args.img); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractDominantColorsKmeans() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPixels(t *testing.T) {
	type args struct {
		img image.Image
	}
	tests := []struct {
		name string
		args args
		want []color.Color
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPixels(tt.args.img); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPixels() = %v, want %v", got, tt.want)
			}
		})
	}
}
