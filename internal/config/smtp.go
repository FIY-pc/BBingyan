package config

type SmtpConfig struct {
	SmtpUser     string             `yaml:"smtpUser"`
	SmtpNickname string             `yaml:"smtpNickname"`
	SmtpPassword string             `yaml:"smtpPassword"`
	SmtpHost     string             `yaml:"smtpHost"`
	SmtpPort     string             `yaml:"smtpPort"`
	Captcha      *CaptchaConfig     `yaml:"captchaConfig"`
	WeeklyEmail  *WeeklyEmailConfig `yaml:"weeklyEmail"`
}

type CaptchaConfig struct {
	Expire   string `yaml:"expire"`
	Interval string `yaml:"interval"`
}

type WeeklyEmailConfig struct {
	RoutineNum int    `yaml:"routineNum"`
	TimeOut    string `yaml:"timeOut"`
	RateLimit  string `yaml:"rateLimit"`
}
