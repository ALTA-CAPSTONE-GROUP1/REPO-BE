package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	JWT                    string
	CloudinaryName         string
	CloudinaryApiKey       string
	CloudinaryApiScret     string
	CloudinaryUploadFolder string
	TokenSuperAdmin        string
	EmailSecret            string
	Email                  string
	PasswordUser           string
	AutoMessageHyperApp    string
)

type AppConfig struct {
	DBUSER     string
	DBPASSWORD string
	DBHOST     string
	DBPORT     string
	DBNAME     string
}

func InitConfig() *AppConfig {
	return ReadEnv()
}

func ReadEnv() *AppConfig {
	app := AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("DBUSER"); found {
		app.DBUSER = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPASSWORD"); found {
		app.DBPASSWORD = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHOST"); found {
		app.DBHOST = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPORT"); found {
		app.DBPORT = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBNAME"); found {
		app.DBNAME = val
		isRead = false
	}

	if val, found := os.LookupEnv("JWT"); found {
		JWT = val
		isRead = false
	}

	if val, found := os.LookupEnv("CLOUDINARY_CLOUD_NAME"); found {
		CloudinaryName = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_API_KEY"); found {
		CloudinaryApiKey = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_API_SECRET"); found {
		CloudinaryApiScret = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_UPLOAD_FOLDER"); found {
		CloudinaryUploadFolder = val
		isRead = false
	}

	if val, found := os.LookupEnv("TOKEN_SUPER_ADMIN"); found {
		TokenSuperAdmin = val
		isRead = false
	}

	if val, found := os.LookupEnv("EMAIL_SECRET"); found {
		TokenSuperAdmin = val
		isRead = false
	}
	if val, found := os.LookupEnv("EMAIL"); found {
		TokenSuperAdmin = val
		isRead = false
	}

	if val, found := os.LookupEnv("PASSWORD_USER"); found {
		PasswordUser = val
		isRead = false
	}

	if val, found := os.LookupEnv("AUTO_MESSAGE"); found {
		AutoMessageHyperApp = val
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error read config : ", err.Error())
			return nil
		}

		app.DBUSER = viper.Get("DBUSER").(string)
		app.DBPASSWORD = viper.Get("DBPASSWORD").(string)
		app.DBHOST = viper.Get("DBHOST").(string)
		app.DBPORT = viper.Get("DBPORT").(string)
		app.DBNAME = viper.Get("DBNAME").(string)

		JWT = viper.Get("JWT").(string)

		CloudinaryName = viper.Get("CLOUDINARY_CLOUD_NAME").(string)
		CloudinaryApiKey = viper.Get("CLOUDINARY_API_KEY").(string)
		CloudinaryApiScret = viper.Get("CLOUDINARY_API_SECRET").(string)
		CloudinaryUploadFolder = viper.Get("CLOUDINARY_UPLOAD_FOLDER").(string)
		TokenSuperAdmin = viper.Get("TOKEN_SUPER_ADMIN").(string)
		EmailSecret = viper.Get("EMAIL_SECRET").(string)
		Email = viper.Get("EMAIL").(string)
		PasswordUser = viper.Get("PASSWORD_USER").(string)
		AutoMessageHyperApp = viper.Get("AUTO_MESSAGE").(string)

	}
	return &app
}
