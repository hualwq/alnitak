package config

type RateLimit struct {
	Enabled  bool  `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Capacity int   `mapstructure:"capacity" json:"capacity" yaml:"capacity"`
	Rate     int   `mapstructure:"rate" json:"rate" yaml:"rate"`
	Timeout  int64 `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}

type Security struct {
	AccessJwtSecret          string     `mapstructure:"access_jwt_secret" json:"access_jwt_secret" yaml:"access_jwt_secret"`
	RefreshJwtSecret         string     `mapstructure:"refresh_jwt_secret" json:"refresh_jwt_secret" yaml:"refresh_jwt_secret"`
	CloseRecordUserOperation bool       `mapstructure:"close_record_user_operation" json:"close_record_user_operation" yaml:"close_record_user_operation"`
	RateLimit                RateLimit  `mapstructure:"rate_limit" json:"rate_limit" yaml:"rate_limit"`
}
