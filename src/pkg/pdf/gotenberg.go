package pdf

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewGotenbergClient(url string) *Client {
	return &Client{BaseURL: url}
}

func (c *Client) HTMLToPDF(htmlContent string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("files", "index.html")
	part.Write([]byte(htmlContent))

	writer.WriteField("paperWidth", "297mm")
	writer.WriteField("paperHeight", "210mm")
	writer.WriteField("marginTop", "0")
	writer.WriteField("marginBottom", "0")
	writer.WriteField("marginLeft", "0")
	writer.WriteField("marginRight", "0")
	writer.WriteField("printBackground", "true")

	writer.Close()

	req, _ := http.NewRequest("POST", c.BaseURL+"forms/chromium/convert/html", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("could not convert html to pdf: %s", string(b))
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
