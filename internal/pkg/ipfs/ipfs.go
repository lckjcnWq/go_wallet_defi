package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type IPFSClient struct {
	endpoint string
	gateway  string
}

func NewIPFSClient(endpoint, gateway string) *IPFSClient {
	return &IPFSClient{endpoint: endpoint, gateway: gateway}
}

// UploadFile 上传文件到IPFS
func (c *IPFSClient) UploadFile(file multipart.File) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/api/v0/add", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Hash string `json:"Hash"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Hash, nil
}

// UploadJSON 上传JSON数据到IPFS
func (c *IPFSClient) UploadJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "metadata.json")
	if err != nil {
		return "", err
	}

	_, err = part.Write(jsonData)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/api/v0/add", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Hash string `json:"Hash"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Hash, nil
}

// GetGatewayURL 获取IPFS网关URL
func (c *IPFSClient) GetGatewayURL(hash string) string {
	return fmt.Sprintf("%s/ipfs/%s", c.gateway, hash)
}
