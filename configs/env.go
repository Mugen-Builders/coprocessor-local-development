package configs

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	AnvilWsURL                string
	DappAddress               string
	AnvilHttpURL              string
	InputBoxAddress           string
	AnvilInputBoxBlock        string
	CoprocessorMachineHash    string
	MockCoprocessorAddress    string
	CoprocessorAdapterAddress string
}

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

func LoadConfig(path string) (*Config, error) {
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
			return nil, err
		}
	}

	envVars := &Config{
		AnvilWsURL:                verifyEnv("ANVIL_WS_URL"),
		DappAddress:               "0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e",
		AnvilHttpURL:              verifyEnv("ANVIL_HTTP_URL"),
		InputBoxAddress:           "0x59b22D57D4f067708AB0c00552767405926dc768",
		AnvilInputBoxBlock:        verifyEnv("ANVIL_INPUT_BOX_BLOCK"),
		CoprocessorMachineHash:    verifyEnv("COPROCESSOR_MACHINE_HASH"),
		CoprocessorAdapterAddress: verifyEnv("COPROCESSOR_ADAPTER_ADDRESS"),
		MockCoprocessorAddress:    "0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE",
	}

	return envVars, nil
}

func verifyEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		slog.Error(fmt.Sprintf("%s environment variable not set or empty", key))
		os.Exit(1)
	}
	return val
}
