package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	imagePath := "D:\\projects\\server\\captcha.png"
	client.SetImage(imagePath)
	captchaText, err := client.Text()
	if err != nil {
		fmt.Printf("Error recognizing captcha: %v\n", err)
		return
	}

	fmt.Printf("Recognized captcha text: %s\n", captchaText)
}
