package response

type GetMatchResponse struct {
	ID         int    `json:"id"`
	HomeTeamID int    `json:"home_team_id"`
	AwayTeamID int    `json:"away_team_id"`
	UTCDate    string `json:"utc_date"`
	Status     string `json:"status"`
}

type MatchResponse struct {
	Matches []ApiMatch `json:"matches"`
}

type ApiMatch struct {
	ID       int     `json:"id"`
	UtcDate  string  `json:"utcDate"`
	Status   string  `json:"status"`
	Matchday int     `json:"matchday"`
	Stage    string  `json:"stage"`
	HomeTeam ApiTeam `json:"homeTeam"`
	AwayTeam ApiTeam `json:"awayTeam"`
	Venue    string  `json:"venue"`
}

type ApiTeam struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	TLA       string `json:"tla"`
}
