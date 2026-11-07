package config

import "testing"

func validConfig() Config {
	return Config{
		App:      AppConfig{Name: "ai-analysis", LogLevel: "info"},
		Database: DatabaseConfig{URLEnv: "DATABASE_URL", ConnectTimeoutSeconds: 5},
		Schedule: ScheduleConfig{EveryHours: 6, Timezone: "Asia/Kolkata", MaxRunMinutes: 20},
		Reports:  ReportsConfig{Directory: "reports"},
	}
}

func TestValidateRejectsDuplicateSourceIDs(t *testing.T) {
	cfg := validConfig()
	cfg.DirectCareerSources = []SourceConfig{{ID: "same"}}
	cfg.ThirdPartySources = []SourceConfig{{ID: "same"}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected duplicate source ID error")
	}
}

func TestValidateAcceptsFoundationConfig(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}
