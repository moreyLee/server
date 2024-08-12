package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	apiURL := "http://localhost:8000/ocr"
	imagePath := "/Users/david/Downloads/projects/server/captcha.png"

	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		panic(err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	data := url.Values{}
	data.Set("image", base64Image)
	data.Set("probability", "false")
	data.Set("png_fix", "false")

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
