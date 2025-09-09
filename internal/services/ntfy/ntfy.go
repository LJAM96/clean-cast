package ntfy

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"

	"ikoyhn/podcast-sponsorblock/internal/config"
)

type ntfyConfigStruct struct {
	Server string `yaml:"server"`
	Topic  string `yaml:"topic"`
}

var ntfyConfig ntfyConfigStruct

func init() {
	configPath := filepath.Join(config.Config.ConfigDir, "properties.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create empty properties.yml if it doesn't exist
		emptyConfig := ntfyConfigStruct{}
		data, err := yaml.Marshal(emptyConfig)
		if err != nil {
			log.Error("Failed to marshal empty config:", err)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			log.Error("Failed to create empty properties.yml:", err)
		}
		ntfyConfig = emptyConfig
		return
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Error("Failed to read properties.yml:", err)
		return
	}
	if err := yaml.Unmarshal(data, &ntfyConfig); err != nil {
		log.Error("Failed to unmarshal properties.yml:", err)
		return
	}
}

func SendNotification(message, title string) error {
	if ntfyConfig.Server == "" || ntfyConfig.Topic == "" {
		return nil
	}
	url := fmt.Sprintf("%s/%s", ntfyConfig.Server, ntfyConfig.Topic)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	if title != "" {
		req.Header.Set("Title", title)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("ntfy returned status: %s", resp.Status)
	}
	return nil
}
