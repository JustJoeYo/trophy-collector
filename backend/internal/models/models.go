package models

type Player struct {
    AccountID   uint32 `json:"account_id"`
    PlayerName  string `json:"player_name"`
    AvatarURL   string `json:"avatar_url"`
    Wins        int    `json:"wins"`
    Losses      int    `json:"losses"`
}

type Match struct {
    MatchID       uint64 `json:"match_id"`
    HeroID        uint32 `json:"hero_id"`
    Outcome       string `json:"outcome"`
    Kills         int    `json:"kills"`
    Deaths        int    `json:"deaths"`
    Assists       int    `json:"assists"`
    DurationSecs  int    `json:"duration_secs"`
    StartTime     int64  `json:"start_time"`
}

type Hero struct {
    HeroID    uint32 `json:"id"`
    ClassName string `json:"class_name"`
    Name      string `json:"name"`
}
