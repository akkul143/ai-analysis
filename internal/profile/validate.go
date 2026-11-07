package profile

import (
	"fmt"
	"strings"
)

func (profile Profile) Validate() error {
	if strings.TrimSpace(profile.Identity.Name) == "" {
		return fmt.Errorf("identity.name is required")
	}
	if len(profile.Titles.Preferred) == 0 {
		return fmt.Errorf("titles.preferred must contain at least one title")
	}
	if len(profile.Locations.Preferred) == 0 && len(profile.Locations.Acceptable) == 0 {
		return fmt.Errorf("at least one preferred or acceptable location is required")
	}
	if profile.Experience.MinimumJobYears < 0 {
		return fmt.Errorf("experience.minimumJobYears cannot be negative")
	}
	if profile.Experience.MaximumJobYears < profile.Experience.MinimumJobYears {
		return fmt.Errorf("experience.maximumJobYears cannot be less than minimumJobYears")
	}
	if len(profile.Skills.Strong) == 0 && len(profile.Skills.Working) == 0 {
		return fmt.Errorf("at least one strong or working skill is required")
	}
	if err := validateUniqueSkills(profile.Skills); err != nil {
		return err
	}
	return nil
}

func validateUniqueSkills(skills Skills) error {
	groups := []struct {
		name   string
		values []string
	}{
		{name: "strong", values: skills.Strong},
		{name: "working", values: skills.Working},
		{name: "learning", values: skills.Learning},
		{name: "notUsed", values: skills.NotUsed},
	}
	seen := make(map[string]string)
	for _, group := range groups {
		for _, value := range group.values {
			key := strings.ToLower(strings.TrimSpace(value))
			if key == "" {
				return fmt.Errorf("skills.%s contains an empty value", group.name)
			}
			if previous, exists := seen[key]; exists {
				return fmt.Errorf("skill %q appears in both %s and %s", value, previous, group.name)
			}
			seen[key] = group.name
		}
	}
	return nil
}
