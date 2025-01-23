package configs

import (
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

func readTOML(name string) string {
	bytes, err := os.ReadFile(name)
	if err != nil {
		slog.Error("Failed to read file", "error", err)
		os.Exit(1)
	}
	return string(bytes)
}

type configTOML = map[string]map[string]string

func decodeTOML(data string) configTOML {
	var config configTOML
	_, err := toml.Decode(data, &config)
	if err != nil {
		slog.Error("Failed to decode TOML", "error", err)
		os.Exit(1)
	}
	return config
}

func sortConfig(config configTOML) []string {
	var keys []string
	for section, variables := range config {
		for key := range variables {
			keys = append(keys, section+"."+key)
		}
	}
	sort.Strings(keys)
	return keys
}

func LoadConfig(path string) error {
	data := readTOML(path)
	config := decodeTOML(data)

	sortedKeys := sortConfig(config)

	for _, key := range sortedKeys {
		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			continue
		}
		section, variable := parts[0], parts[1]

		value := config[section][variable]
		envName := strings.ToUpper(section + "_" + variable)
		err := os.Setenv(envName, value)
		if err != nil {
			return err
		}
	}
	return nil
}
