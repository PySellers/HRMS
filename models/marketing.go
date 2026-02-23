package models

// MarketingActivity represents a single daily outreach or procurement action
type MarketingActivity struct {
	ID               int     `json:"id"`
	Date             string  `json:"date"`
	PersonName       string  `json:"person_name"`
	Role             string  `json:"role"`
	OutreachSource   string  `json:"outreach_source"`
	PlatformName     string  `json:"platform_name"`
	ActivityType     string  `json:"activity_type"`
	Industry         string  `json:"industry"`
	ProjectType      string  `json:"project_type"`
	EstimatedValue   float64 `json:"estimated_value"`
	URL              string  `json:"url"`
	Status           string  `json:"status"`
	NextFollowUpDate string  `json:"next_followup_date"`
	KeyNotes         string  `json:"key_notes"`
	TimeSpent        int     `json:"time_spent"`
}

// MarketingDB holds all marketing data (separate JSON file)
type MarketingDB struct {
	Activities []MarketingActivity `json:"activities"`
}
