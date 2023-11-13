package config

type IConfig interface {
	GetPlugins() []PluginConfig
	FindPlugin(name string) *PluginConfig
}

type PluginConfig struct {
	Type   string
	Name   string
	Config map[string]string
}

type Config struct {
	Plugins []PluginConfig
}

func (c *Config) GetPlugins() []PluginConfig {
	return c.Plugins
}

func (c *Config) FindPlugin(name string) *PluginConfig {
	for _, p := range c.Plugins {
		if p.Name == name {
			return &p
		}
	}
	return nil
}
