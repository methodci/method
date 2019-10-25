package method

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadRequest struct {
	Title string `json:"name"`

	PWD string `json:"pwd"`
	SHA string `json:"sha"`

	Time time.Time `json:"time"`
}

func Send(endpoint string, details UploadRequest, files ...string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, f := range files {
		f1, err := writer.CreateFormFile("file[]", filepath.Base(f))
		if err != nil {
			return err
		}

		of, err := os.Open(f)
		defer of.Close()
		if err != nil {
			return err
		}

		io.Copy(f1, of)
	}

	detailstr, err := json.Marshal(details)
	if err != nil {
		return err
	}

	// writer.WriteField("file2", "fresh")
	writer.WriteField("data", string(detailstr))
	writer.Close()

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return err
	}

	// Headers
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("analysis error: %s", respBody)
	}

	fmt.Println(string(respBody))

	return nil
}
