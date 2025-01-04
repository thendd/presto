package discord

type Guild struct {
	ID                          string        `json:"id"`
	Name                        string        `json:"name"`
	Icon                        string        `json:"icon"`
	IconHash                    string        `json:"icon_hash"`
	Splash                      string        `json:"splash"`
	DiscoverySplash             string        `json:"discovery_splash"`
	Owner                       bool          `json:"owner"`
	OwnerID                     string        `json:"owner_id"`
	Permissions                 string        `json:"permissions"`
	Region                      string        `json:"region"`
	AfkChannelID                string        `json:"afk_channel_id"`
	AfkTimeout                  int           `json:"afk_timeout"`
	WidgetEnabled               bool          `json:"widget_enabled"`
	WidgetChannelID             string        `json:"widget_channel_id"`
	VerificationLevel           int           `json:"verification_level"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Roles                       []Role        `json:"roles"`
	Emojis                      []Emoji       `json:"emojis"`
	Features                    []string      `json:"features"`
	MfaLevel                    int           `json:"mfa_level"`
	ApplicationID               string        `json:"application_id"`
	SystemChannelID             string        `json:"system_channel_id"`
	SystemChannelFlags          int           `json:"system_channel_flags"`
	RulesChannelID              string        `json:"rules_channel_id"`
	JoinedAt                    string        `json:"joined_at"`
	Large                       bool          `json:"large"`
	Unavailable                 bool          `json:"unavailable"`
	MemberCount                 int           `json:"member_count"`
	Members                     []GuildMember `json:"members"`
	Channels                    []Channel     `json:"channels"`
	Threads                     []Channel     `json:"threads"`
	MaxPresences                int           `json:"max_presences"`
	MaxMembers                  int           `json:"max_members"`
	VanityURLCode               string        `json:"vanity_url_code"`
	Description                 string        `json:"description"`
	Banner                      string        `json:"banner"`
	PremiumTier                 int           `json:"premium_tier"`
	PremiumSubscriptionCount    int           `json:"premium_subscription_count"`
	PreferredLocale             string        `json:"preferred_locale"`
	PublicUpdatesChannelID      string        `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int           `json:"max_video_channel_users"`
	ApproximateMemberCount      int           `json:"approximate_member_count"`
	ApproximatePresenceCount    int           `json:"approximate_presence_count"`
	NSFWLevel                   int           `json:"nsfw_level"`
	Stickers                    []Sticker     `json:"stickers"`
	BoostProgressBarEnabled     bool          `json:"premium_progress_bar_enabled"`
}

type Sticker struct {
	ID          string `json:"id"`
	PackID      string `json:"pack_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Type        int    `json:"type"`
	FormatType  int    `json:"format_type"`
	Available   bool   `json:"available"`
	GuildID     string `json:"guild_id"`
	User        User   `json:"user"`
	SortValue   int    `json:"sort_value"`
}
