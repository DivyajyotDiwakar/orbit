package configs

import (
	"log"
	"os"
)

type PortConfig struct {
	Value string
}

type MongoDbConfig struct {
	Url string
}

type JwtConfig struct {
	Key string
}

type R2Config struct {
	Bucket     string
	Access_Key string
	Secret_Key string
	Endpoint   string
	Domain     string
}

type Redis struct {
	Url      string
	// Password string
}

type Config struct {
	PORT    PortConfig
	Mongo   MongoDbConfig
	JWT_KEY JwtConfig
	R2      R2Config
	Redis   Redis
}

func LoadConfig() *Config {
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8000"
	}
	Port := PortConfig{Value: portString}

	MongoDbUrlString := os.Getenv("MONGODBURI")
	if MongoDbUrlString == "" {
		log.Fatal("Doesn't find any database url to connect with..!!")
	}
	MongoDb := MongoDbConfig{
		Url: MongoDbUrlString,
	}

	jwtString := os.Getenv("JWT_KEY")
	if jwtString == "" {
		log.Fatal("Couldn't Find jwt key")
	}
	Jwt := JwtConfig{
		Key: jwtString,
	}

	R2_ACCESS_Key := os.Getenv("R2_ACCESS_KEY_ID")
	R2_SECRET_ACCESS_KEY := os.Getenv("R2_SECRET_ACCESS_KEY")
	R2_ENDPOINT := os.Getenv("R2_ENDPOINT")
	R2_BUCKET_NAME := os.Getenv("R2_BUCKET_NAME")
	R2_DOMAIN := os.Getenv("R2_DOMAIN")
	if R2_ACCESS_Key == "" || R2_SECRET_ACCESS_KEY == "" || R2_ENDPOINT == "" || R2_BUCKET_NAME == "" {
		log.Fatal("R2 configuration is incomplete. Please set R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, and R2_ENDPOINT environment variables.")
	}
	R2 := R2Config{
		Access_Key: R2_ACCESS_Key,
		Secret_Key: R2_SECRET_ACCESS_KEY,
		Endpoint:   R2_ENDPOINT,
		Bucket:     R2_BUCKET_NAME,
		Domain:     R2_DOMAIN,
	}

	redis_url := os.Getenv("REDIS_URL")
	// redis_pass := os.Getenv("REDIS_PASS")
	// if redis_url == ""  || redis_pass == "" {
	// 	log.Fatal("No Redis Url Found..!!")
	// }
	redis := Redis{
		Url:      redis_url,
		// Password: redis_pass,
	}

	return &Config{
		PORT:    Port,
		Mongo:   MongoDb,
		JWT_KEY: Jwt,
		R2:      R2,
		Redis:   redis,
	}
}
