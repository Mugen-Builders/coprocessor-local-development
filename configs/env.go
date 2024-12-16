package configs

import (
	"github.com/spf13/viper"
)

type CFG struct {
	GRAPHQL_URL					 string `mapstructure:"GRAPHQL_URL"`
	RPC_URL_WS                      string `mapstructure:"RPC_URL_WS"`
	PRIVATE_KEY                     string `mapstructure:"PRIVATE_KEY"`
	RPC_URL_HTTP                    string `mapstructure:"RPC_URL_HTTP"`
	INPUT_BOX_ADDRESS               string `mapstructure:"INPUT_BOX_ADDRESS"`
	COPROCESSOR_CALLER_MOCK_ADDRESS string `mapstructure:"COPROCESSOR_CALLER_MOCK_ADDRESS"`
}

func LoadConfig() (*CFG, error) {
	var cfg *CFG

	viper.SetDefault("GRAPHQL_URL", "http://127.0.0.1:8080/graphql")
	viper.SetDefault("RPC_URL_WS", "ws://127.0.0.1:8545")
	viper.SetDefault("RPC_URL_HTTP", "http://127.0.0.1:8545")
	viper.SetDefault("INPUT_BOX_ADDRESS", "0x59b22D57D4f067708AB0c00552767405926dc768")
	viper.SetDefault("PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")

	viper.SetConfigName("app_config")
	viper.AutomaticEnv()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
