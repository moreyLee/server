package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tebeka/selenium"
	"log"
	"net/http"
	"time"
)

const (
	seleniumPath     = "path/to/selenium-server-standalone.jar" // Path to Selenium Server JAR file
	chromeDriverPath = "path/to/chromedriver"                   // Path to ChromeDriver
	sport            = 8080
)

func main() {
	r := gin.Default()

	// Define a route to handle Selenium actions
	r.GET("/selenium", func(c *gin.Context) {
		if err := runSelenium(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Selenium action executed successfully"})
	})

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func runSelenium() error {
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),             // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(chromeDriverPath), // Specify the path to ChromeDriver in order to use Chrome.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, sport, opts...)
	if err != nil {
		return fmt.Errorf("error starting the Selenium server: %v", err)
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", sport))
	if err != nil {
		return fmt.Errorf("error creating new WebDriver session: %v", err)
	}
	defer wd.Quit()

	// Example Selenium action: Navigate to a website and print the title
	if err := wd.Get("https://www.godaddy.com/"); err != nil {
		return fmt.Errorf("error navigating to example.com: %v", err)
	}

	time.Sleep(2 * time.Second) // Wait for the page to load

	title, err := wd.Title()
	if err != nil {
		return fmt.Errorf("error getting page title: %v", err)
	}
	log.Printf("Page title: %s", title)

	return nil
}
