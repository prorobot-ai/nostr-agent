package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// BotConfig defines the structure for each bot
type BotConfig struct {
	Name          string        `yaml:"name"`
	Aliases       []string      `yaml:"aliases"`
	RelayURL      string        `yaml:"relay_url"`
	Nsec          string        `yaml:"nsec"`
	ChannelID     string        `yaml:"channel_id"`
	Listener      string        `yaml:"listener"`
	Publisher     string        `yaml:"publisher"`
	Handler       string        `yaml:"handler"`
	EventType     string        `yaml:"event_type"`
	ProgramConfig ProgramConfig `yaml:"program"`
}

// BotConfigs is a wrapper to handle multiple bots
type BotConfigs struct {
	Bots []BotConfig `yaml:"bots"`
}

type ProgramConfig struct {
	MaxRunCount   int          `yaml:"max_run_count"`
	ResponseDelay int          `yaml:"response_delay"`
	WorkerConfig  WorkerConfig `yaml:"worker"`
	HubConfig     HubConfig    `yaml:"hub"`
}

type HubConfig struct {
	Socket string `yaml:"socket"`
}

type WorkerConfig struct {
	Url     string `yaml:"url"`
	Address string `yaml:"address"`
}

// LoadBotConfigs loads the bot configurations from a YAML file
func LoadBotConfigs(path string) (*BotConfigs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
		return nil, err
	}

	var botConfigs BotConfigs
	if err := yaml.Unmarshal(data, &botConfigs); err != nil {
		log.Fatalf("❌ Failed to parse YAML config: %v", err)
		return nil, err
	}

	return &botConfigs, nil
}
