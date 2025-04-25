package setting

type Config struct {
	Server          ServerSetting         `mapstructure:"server"`
	Postgres        PostgresSetting       `mapstructure:"postgres"`
	PostgresSetting PostgresSetting       `mapstructure:"postgres_setting"`
	Logger          LogSetting            `mapstructure:"logger"`
	Security        SecuritySetting       `mapstructure:"security"`
	RabbitMQ        RabbitMQConfig        `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	Cronjob         CronjobSetting        `mapstructure:"cronjob"`
	TokenValidation TokenValidationConfig `mapstructure:"token_validation"`
	Redis           RedisSetting          `mapstructure:"redis"`
	Exchange        ExchangeSetting       `mapstructure:"exchange"`
	Queue           QueueSetting          `mapstructure:"queue"`
}

type RedisSetting struct {
	URL      string `mapstructure:"url"`
	PoolSize int    `mapstructure:"pool_size"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type ServerSetting struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	Mode         string `mapstructure:"mode"`
}

type PostgresSetting struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"db_name"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	TimeZone        string `mapstructure:"time_zone"`
}

type LogSetting struct {
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type SecuritySetting struct {
	CryptoKeys CryptoKeysSetting `mapstructure:"crypto_keys"`
}

type CryptoKeysSetting struct {
	Asymmetric AsymmetricSetting `mapstructure:"asymmetric"`
	Symmetric  SymmetricSetting  `mapstructure:"symmetric"`
}

type AsymmetricSetting struct {
	PrivKey      string `mapstructure:"priv_key"`
	SenderPubKey string `mapstructure:"sender_pub_key"`
	PubKey       string `mapstructure:"pub_key"`
}

type SymmetricSetting struct {
	AESKey string `mapstructure:"aes_key"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	VHost    string `mapstructure:"vhost" json:"vhost" yaml:"vhost"`
}

type CronjobSetting struct {
	CronExecuteInterest string `mapstructure:"cron_execute_interest"`
}

type TokenValidationConfig struct {
	ClientID string `mapstructure:"client_id"`
	URL      string `mapstructure:"url"`
	XApiKey  string `mapstructure:"x_api_key"`
}

type ExchangeSetting struct {
	Test string `mapstructure:"test"`
}

type QueueSetting struct {
	Test string `mapstructure:"test"`
}
