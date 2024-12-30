package config

import (
	"github.com/G-Villarinho/food-shop-api/config/models"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

var Env models.Environment

func LoadEnvironments() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&Env)
	if err != nil {
		panic(err)
	}

}
