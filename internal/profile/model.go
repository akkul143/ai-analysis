package profile

type Profile struct {
	Identity          Identity          `yaml:"identity"`
	Experience        Experience        `yaml:"experience"`
	Titles            Titles            `yaml:"titles"`
	Locations         Locations         `yaml:"locations"`
	Skills            Skills            `yaml:"skills"`
	Domains           Domains           `yaml:"domains"`
	Employment        Employment        `yaml:"employment"`
	WorkAuthorization WorkAuthorization `yaml:"workAuthorization"`
}

type Identity struct {
	Name string `yaml:"name"`
}

type Experience struct {
	MinimumJobYears float64 `yaml:"minimumJobYears"`
	MaximumJobYears float64 `yaml:"maximumJobYears"`
}

type Titles struct {
	Preferred  []string `yaml:"preferred"`
	Acceptable []string `yaml:"acceptable"`
	Excluded   []string `yaml:"excluded"`
}

type Locations struct {
	Preferred []string `yaml:"preferred"`
	Acceptable []string `yaml:"acceptable"`
	Remote    bool     `yaml:"remote"`
	Hybrid    bool     `yaml:"hybrid"`
	Onsite    bool     `yaml:"onsite"`
	Relocation bool    `yaml:"relocation"`
}

type Skills struct {
	Strong   []string `yaml:"strong"`
	Working  []string `yaml:"working"`
	Learning []string `yaml:"learning"`
	NotUsed  []string `yaml:"notUsed"`
}

type Domains struct {
	Preferred []string `yaml:"preferred"`
}

type Employment struct {
	Types []string `yaml:"types"`
}

type WorkAuthorization struct {
	Countries []string `yaml:"countries"`
}
