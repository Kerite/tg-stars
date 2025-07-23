package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ExportMemoryBody struct {
	UserId string `json:"user_id"`
}

// Download file from the given URL and return the file path.
func SaveFile(url string) (string, error) {
	return "", nil
}

func GetFileData(downloadURL string) ([]byte, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return data, nil
}

func ExportMemory(userID string) ([]byte, error) {
	jsonBody, err := json.Marshal(&ExportMemoryBody{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	// Export memory from the backend service
	httpReq, err := http.NewRequest(http.MethodPost, os.Getenv("BACKEND_URL")+"/api/export-memory", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	os.WriteFile("exported_memory_"+userID+".snapshot", data, 0644)
	if httpRsp.StatusCode != http.StatusOK {
		fmt.Println("Error response from export memory:", httpRsp.Status)
		return nil, err
	}
	return data, nil
}

func ImportMemory(userID string, fileData []byte) error {
	bodyBuffer := bytes.NewBuffer([]byte{})

	mpWriter := multipart.NewWriter(bodyBuffer)
	mpWriter.WriteField("user_id", userID)

	part, err := mpWriter.CreateFormFile("snapshot", "imported_memory_"+userID+".snapshot")
	if err != nil {
		return err
	}
	io.Copy(part, bytes.NewReader(fileData))

	defer mpWriter.Close()

	reqBody, err := http.NewRequest(http.MethodPost, os.Getenv("BACKEND_URL")+"/api/import-memory", bodyBuffer)
	if err != nil {
		return err
	}
	reqBody.Header.Set("Content-Type", mpWriter.FormDataContentType())

	httpRsp, err := http.DefaultClient.Do(reqBody)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()

	data, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if httpRsp.StatusCode != http.StatusOK {
		fmt.Println("Error response from import memory:", httpRsp.Status, " - ", string(data))
		return fmt.Errorf("failed to import memory: %s", httpRsp.Status)
	}
	return nil
}

func Chat(userId string, message string) (string, error) {
	jsonBody, err := json.Marshal(map[string]string{
		"user_id": userId,
		"message": message,
	})
	if err != nil {
		return "", err
	}
	reqBody, err := http.NewRequest(http.MethodPost, os.Getenv("BACKEND_URL")+"/api/chat", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return "", err
	}
	reqBody.Header.Set("Content-Type", "application/json")
	httpRsp, err := http.DefaultClient.Do(reqBody)
	if err != nil {
		return "", err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		fmt.Println("Error response from send message:", httpRsp.Status)
		return "", fmt.Errorf("failed to send message: %s", httpRsp.Status)
	}
	data, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	fmt.Println("Response data:", string(data))
	var response map[string]string
	if err := json.Unmarshal(data, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	fmt.Println(response)
	if msg, ok := response["message"]; ok {
		return msg, nil
	}
	return  "", nil
}
