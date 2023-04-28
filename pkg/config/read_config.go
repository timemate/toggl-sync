package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig() (IConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	//viper.AddConfigPath("$HOME/.toggl-sync")
	//viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	plugins := viper.Get("plugins").([]interface{})
	cfg := &Config{
		Plugins: make([]PluginConfig, len(plugins), len(plugins)),
	}
	fmt.Printf("%v", plugins)
	for i, p := range plugins {
		plugin := p.(map[interface{}]interface{})
		config := PluginConfig{
			Type:   plugin["type"].(string),
			Name:   plugin["name"].(string),
			Config: make(map[string]string, 0),
		}
		params := plugin["config"].(map[interface{}]interface{})
		for k, v := range params {
			config.Config[k.(string)] = v.(string)
		}
		cfg.Plugins[i] = config
	}

	return cfg, err
}
