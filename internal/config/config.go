package config

type Config struct {
	App                  AppConfig      `yaml:"app"`
	Database             DatabaseConfig `yaml:"database"`
	Schedule             ScheduleConfig `yaml:"schedule"`
	BrowserWorker        BrowserConfig  `yaml:"browserWorker"`
	Reports              ReportsConfig  `yaml:"reports"`
	DirectCareerSources  []SourceConfig `yaml:"directCareerSources"`
	ThirdPartySources    []SourceConfig `yaml:"thirdPartySources"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"logLevel"`
}

type DatabaseConfig struct {
	URLEnv                string `yaml:"urlEnv"`
	ConnectTimeoutSeconds int    `yaml:"connectTimeoutSeconds"`
}

type ScheduleConfig struct {
	EveryHours    int    `yaml:"everyHours"`
	Timezone      string `yaml:"timezone"`
	MaxRunMinutes int    `yaml:"maxRunMinutes"`
}

type BrowserConfig struct {
	Enabled              bool     `yaml:"enabled"`
	Command              string   `yaml:"command"`
	Args                 []string `yaml:"args"`
	HealthTimeoutSeconds int      `yaml:"healthTimeoutSeconds"`
}

type ReportsConfig struct {
	Directory string   `yaml:"directory"`
	Formats   []string `yaml:"formats"`
}

type SourceConfig struct {
	ID         string `yaml:"id"`
	Company    string `yaml:"company"`
	Adapter    string `yaml:"adapter"`
	CareersURL string `yaml:"careersUrl"`
	Enabled    bool   `yaml:"enabled"`
}
