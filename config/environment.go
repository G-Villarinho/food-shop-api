package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

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

	Env.PrivateKey, err = loadPrivateKey()
	if err != nil {
		panic(err)
	}

	Env.PublicKey, err = loadPublicKey()
	if err != nil {
		panic(err)
	}
}

func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile("ec_private_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("erro ao decodificar o bloco PEM contendo a chave privada")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadPublicKey() (*ecdsa.PublicKey, error) {
	keyData, err := os.ReadFile("ec_public_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("erro ao decodificar o bloco PEM contendo a chave pública")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("não é uma chave pública ECDSA")
	}

	return ecdsaPubKey, nil
}
