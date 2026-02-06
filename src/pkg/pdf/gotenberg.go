package pdf

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NEWGotenbergClient(url string) *Client {
	return &Client{BaseURL: url}
}

func (c *Client) HTMLToPDF(htmlContent string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("files", "index.html")
	part.Write([]byte(htmlContent))
	writer.Close()

	req, _ := http.NewRequest("POST", c.BaseURL+"forms/chromium/convert/html", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
