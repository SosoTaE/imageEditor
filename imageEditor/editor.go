package editor

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"media/types"
	"os"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// Returns absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Returns the maximum value of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Returns the minimum value of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func addEmptySpace(img image.Image, emptySize image.Point, position image.Point) image.Image {
	empty := image.NewNRGBA(image.Rect(0, 0, emptySize.X, emptySize.Y))
	black := color.NRGBA{0, 0, 0, 255}
	for x := 0; x < emptySize.X; x++ {
		for y := 0; y < emptySize.Y; y++ {
			empty.Set(x, y, black)
		}
	}
	draw.Draw(empty, img.Bounds().Add(image.Pt(position.X, position.Y)), img, image.Point{}, draw.Over)
	return empty
}

func crop(img image.Image, coordinates [4]int) image.Image {
	x, y, x1, y1 := coordinates[0], coordinates[1], coordinates[2], coordinates[3]
	width, height := abs(x-x1), abs(y-y1)
	if img == nil {
		fmt.Println("Image is nil!")
		panic("img is nil") // or handle appropriately
	}
	x_c := max(width, img.Bounds().Dx())
	y_c := max(height, img.Bounds().Dy())
	emptyImageSize := image.Point{X: x_c, Y: y_c}
	img = addEmptySpace(img, emptyImageSize, image.Point{X: abs(min(x, 0)), Y: abs(min(y, 0))})
	return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
}

func imgToBase64(url string) string {
	file, err := os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Step 2: Read the file contents
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	return bytesToBase64(fileBytes)
}

func bytesToBase64(bytesData []byte) string {
	return base64.StdEncoding.EncodeToString(bytesData)
}

func base64ToBytes(base64string string) []byte {
	data, _ := base64.StdEncoding.DecodeString(base64string)
	return data
}

func Base64ToImage(base64String string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ImageToBase64(img image.Image) (string, error) {
	buffer := new(bytes.Buffer)

	// Encode the image to JPEG and write to buffer
	err := jpeg.Encode(buffer, img, nil)
	if err != nil {
		return "", err
	}

	// Convert the buffer bytes to a Base64 string
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func ResizeImage(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func ResizeImageInterface(data interface{}, body map[string]interface{}) image.Image {
	var img image.Image
	switch v := data.(type) {
	case image.Image:
		img = v
	case string:
		img, _ = Base64ToImage(v)
	default:
		fmt.Println("Unsupported type:", v)
	}

	size := body["size"].([]interface{})
	width := uint(size[0].(float64))
	height := uint(size[1].(float64))


	return ResizeImage(img, width, height)
}



func cropFunctionInterface(data interface{}, body map[string]interface{}) image.Image {
	var img image.Image
	switch v := data.(type) {
	case image.Image:
		img = v
	case string:
		img, _ = Base64ToImage(v)
	default:
		fmt.Println("Unsupported type:", v)
	}
	
	
	coordinateInterface := body["coordinate"].([]interface{})

	var coordinate [4]int

	for index, each := range coordinateInterface {
		coordinate[index] = int(each.(float64))
	}

	return crop(img, coordinate)
}


var toolDict = map[string]func(interface{}, map[string]interface{}) image.Image{
	// "zoom": nil,
	"crop": cropFunctionInterface,
	"resize": ResizeImageInterface,
}

func imageProcessor(body types.ObjectExample) []string {
	var base64ProcessedImages []string
	for _, base64String := range body.Images {
		var result interface{}
		result = base64String
		for _, each := range body.Function {
			if _, ok := toolDict[each.Name]; !ok {
				fmt.Println("this function does not exist")
			}

			result = toolDict[each.Name](result, each.Parameters)

		}

		data, _ := ImageToBase64(result.(image.Image))
		base64ProcessedImages = append(base64ProcessedImages, data)
	}

	return base64ProcessedImages
}

func Edit(body types.JsonExample) ([][]string, error) {
	objects := body.Objects
	if len(objects) == 0 {
		return nil, errors.New("objects array is empty")
	}

	var base64ImgArray [][]string

	for _, object := range objects {

		if len(object.Images) == 0 {
			return nil, errors.New("images array is empty")
		}

		if len(object.Function) == 0 {
			return nil, errors.New("functions array is empty")
		}

		result := imageProcessor(object)

		base64ImgArray = append(base64ImgArray, result)
	}

	return base64ImgArray, nil
}

// func main() {
// 	start := time.Now()
// 	folderPath := "D:\\pythonProject2\\test_folder"
// 	outputFolderPath := "D:\\pythonProject2\\test_output"
// 	files, _ := ioutil.ReadDir(folderPath)
// 	var urls []string
// 	for _, file := range files {
// 		urls = append(urls, fmt.Sprintf("%s/%s", folderPath, file.Name()))
// 	}

// 	var base64ImagesArray []string
// 	for _, url := range urls {
// 		base64ImagesArray = append(base64ImagesArray, imgToBase64(url))
// 	}

// 	body := map[string]interface{}{
// 		"images": base64ImagesArray,
// 		"crop":   [4]int{-1000, -1000, 2000, 2000},
// 	}

// 	result := edit(body)

// 	for index, each := range result {
// 		bytesData := base64ToBytes(each)
// 		outputPath := fmt.Sprintf("%s/%d.jpg", outputFolderPath, index)
// 		ioutil.WriteFile(outputPath, bytesData, 0644)
// 	}

// 	end := time.Since(start)

// 	fmt.Println(end)

// 	fmt.Println("Press 'Enter' to exit...")
//     fmt.Scanln()  // Wait for user to press 'Enter'
// }

// var base64ImgArray []string

// for imgStr := range imagesInterface {

// 	imgBytes := base64ToBytes(imgStr.(string))
// 	img, _, err := image.Decode(bytes.NewReader(imgBytes))
// 	if err != nil {
// 		fmt.Printf("Error decoding image: %v\n", err)
// 		continue // or handle appropriately
// 	}

// 	for key, valueInterface := range body {
// 		if key == "images" {
// 			continue
// 		}

// 		processor, exists := toolDict[key]
// 		if !exists || processor == nil {
// 			continue
// 		}

// 		if !ok {
// 			continue
// 		}

// 		var value [4]int
// 		tmp := valueInterface.([]interface{})

// 		for i := 0; i < 4; i++ {
// 			floatVal, ok := tmp[i].(float64)
// 			if !ok {
// 				// handle error
// 				continue
// 			}
// 			value[i] = int(floatVal)
// 		}

// 		img = processor(img, value)
// 	}
