package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type Message struct {
	ChatId          string `json:"chatId"`
	Message         string `json:"message"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
}

type FileMessage struct {
	ChatId   string `json:"chatId"`
	UrlFile  string `json:"urlFile"`
	FileName string `json:"fileName"`
}

type TemplateData struct {
	Response string
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.POST("/submit", func(c echo.Context) error {
		idInstance := c.FormValue("idInstance")
		ApiTokenInstance := c.FormValue("ApiTokenInstance")
		action := c.FormValue("action")

		fmt.Println("idInstance: ", idInstance)
		fmt.Println("ApiTokenInstance: ", ApiTokenInstance)

		switch action {
		case "getSettings":
			url := fmt.Sprintf("https://api.green-api.com/waInstance%s/getSettings/%s", idInstance, ApiTokenInstance)

			fmt.Println("url: ", url)

			response, err := http.Get(url)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			fmt.Println("Response data from getSettings: ", string(responseData))

			htmlContent, err := os.ReadFile("static/index.html")
			if err != nil {
				return err
			}

			updatedHtml := strings.Replace(string(htmlContent), "{{.Response}}", (string(responseData)), 1)

			return c.HTMLBlob(http.StatusOK, []byte(updatedHtml))

		case "getStateInstance":
			url := fmt.Sprintf("https://api.green-api.com/waInstance%s/getStateInstance/%s", idInstance, ApiTokenInstance)

			fmt.Println("url: ", url)

			response, err := http.Get(url)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			fmt.Println("Response data from getStateInstance: ", string(responseData))

			htmlContent, err := os.ReadFile("static/index.html")
			if err != nil {
				return err
			}

			updatedHtml := strings.Replace(string(htmlContent), "{{.Response}}", (string(responseData)), 1)

			return c.HTMLBlob(http.StatusOK, []byte(updatedHtml))

		case "sendMessage":
			url := fmt.Sprintf("https://api.green-api.com/waInstance%s/sendMessage/%s", idInstance, ApiTokenInstance)
			number := c.FormValue("number")
			message := c.FormValue("message")

			fmt.Println("url: ", url)
			fmt.Println("message: ", message)
			fmt.Println("number: ", number)

			requestBody := Message{
				ChatId:          number + "@c.us",
				Message:         message,
				QuotedMessageId: "",
			}

			jsonValue, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}

			response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				return err
			}
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			fmt.Println("Response data from sendMessage: ", string(responseData))

			htmlContent, err := os.ReadFile("static/index.html")
			if err != nil {
				return err
			}

			updatedHtml := strings.Replace(string(htmlContent), "{{.Response}}", (string(responseData)), 1)

			return c.HTMLBlob(http.StatusOK, []byte(updatedHtml))

		case "sendFileByURL":
			url := fmt.Sprintf("https://api.green-api.com/waInstance%s/sendFileByUrl/%s", idInstance, ApiTokenInstance)
			number := c.FormValue("numberF")
			fileUrl := c.FormValue("messageF")

			fmt.Println("url: ", url)
			fmt.Println("fileUrl: ", fileUrl)
			fmt.Println("number: ", number)

			requestBody := FileMessage{
				ChatId:   number + "@c.us",
				UrlFile:  fileUrl,
				FileName: filepath.Base(fileUrl),
			}

			fmt.Println("Parsed file name: ", requestBody.FileName)

			jsonValue, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}

			response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				return err
			}
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			fmt.Println("Response data from sendFileByUrl: ", string(responseData))

			htmlContent, err := os.ReadFile("static/index.html")
			if err != nil {
				return err
			}

			updatedHtml := strings.Replace(string(htmlContent), "{{.Response}}", (string(responseData)), 1)

			return c.HTMLBlob(http.StatusOK, []byte(updatedHtml))

		}
		return nil
	})

	e.Logger.Fatal(e.Start("localhost:5000"))

}
