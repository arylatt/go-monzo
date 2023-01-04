package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	CacheFileToken        = "token"
	CacheFileTransactions = "transactions"
)

func LoadCache(fileName string, out any) (err error) {
	filePath := path.Join(viper.GetString("home-dir"), fmt.Sprintf("%s.json", fileName))

	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&out)

	return
}

func SaveCache(fileName string, in any) (err error) {
	filePath := path.Join(viper.GetString("home-dir"), fmt.Sprintf("%s.json", fileName))

	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0700)
}
