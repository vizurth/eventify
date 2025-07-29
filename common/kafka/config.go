package kafka

type KafkaConfig struct {
	Host    string   `yaml:"host" env:"HOST" env-default:"kafka"`
	Port    uint16   `yaml:"port" env:"PORT" env-default:"9092"`
	Brokers []string `yaml:"brokers" env:"BROKERS" env-separator:","`
}

//func NewKafkaConfig() (KafkaConfig, error) {
//	var cfg KafkaConfig
//
//	if err := cleanenv.ReadConfig("configs/configs.yaml", &cfg); err != nil {
//		fmt.Println(err)
//		if err := cleanenv.ReadEnv(&cfg); err != nil {
//			return KafkaConfig{}, fmt.Errorf("error reading Kafka configs: %w", err)
//		}
//	}
//
//	return cfg, nil
//}
