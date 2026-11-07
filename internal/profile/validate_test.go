package profile

import "testing"

func TestValidateRejectsSkillInMultipleGroups(t *testing.T) {
	candidate := Profile{
		Identity:   Identity{Name: "Akkul Gautam"},
		Experience: Experience{MinimumJobYears: 1, MaximumJobYears: 4},
		Titles:     Titles{Preferred: []string{"Software Engineer"}},
		Locations:  Locations{Preferred: []string{"Bengaluru"}},
		Skills: Skills{
			Strong:  []string{"Kafka"},
			Working: []string{"kafka"},
		},
	}

	if err := candidate.Validate(); err == nil {
		t.Fatal("expected duplicate skill validation error")
	}
}

func TestValidateAcceptsMinimumProfile(t *testing.T) {
	candidate := Profile{
		Identity:   Identity{Name: "Akkul Gautam"},
		Experience: Experience{MinimumJobYears: 1, MaximumJobYears: 4},
		Titles:     Titles{Preferred: []string{"Software Engineer"}},
		Locations:  Locations{Preferred: []string{"Bengaluru"}},
		Skills:     Skills{Strong: []string{"Kafka"}},
	}

	if err := candidate.Validate(); err != nil {
		t.Fatalf("expected valid profile, got %v", err)
	}
}
