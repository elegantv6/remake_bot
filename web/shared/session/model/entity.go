package model

type DiscordOAuthSession struct {
	Token string      `json:"token"`
	User  DiscordUser `json:"user"`
}

type DiscordUser struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	GlobalName       string `json:"global_name"`
	DisplayName      string `json:"display_name"`
	Avatar           string `json:"avatar"`
	AvatarDecoration string `json:"avatar_decoration"`
	Discriminator    string `json:"discriminator"`
	PublicFlags      int    `json:"public_flags"`
	Flags            int    `json:"flags"`
	Banner           string `json:"banner"`
	BannerColor      string `json:"banner_color"`
	AccentColor      int    `json:"accent_color"`
	Locale           string `json:"locale"`
	MfaEnabled       bool   `json:"mfa_enabled"`
	PremiumType      int    `json:"premium_type"`
	Email            string `json:"email"`
	Verified         bool   `json:"verified"`
	Bio              string `json:"bio"`
}

type LineOAuthSession struct {
	Token          string   `json:"token"`
	DiscordGuildID string   `json:"discord_guild_id"`
	User           LineUser `json:"user"`
}

type LineToken struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type LineUser struct {
	Iss      string   `json:"iss"`
	Sub      string   `json:"sub"`
	Aud      string   `json:"aud"`
	Exp      int      `json:"exp"`
	Iat      int      `json:"iat"`
	AuthTime int      `json:"auth_time"`
	Nonce    int      `json:"nonce"`
	Amr      []string `json:"amr"`
	Name     string   `json:"name"`
	Picture  string   `json:"picture"`
	Email    string   `json:"email"`
}
