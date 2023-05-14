package domain

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"errors"
	"go.uber.org/zap"
	"os"
	"time"
)

type Ping struct {
	Error bool   `json:"error" bson:"error"`
	Msg   string `json:"msg" bson:"msg"`
}

type PingRequest struct {
	Message string `json:"message"`
}

type Config struct {
	DatabaseUrl  string `mapstructure:"DATABASE_URL"`
	DatabaseName string `mapstructure:"DB_NAME"`
	RedisUrl     string `mapstructure:"REDIS_URL"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func HandleError(c *fiber.Ctx, err error) error {
	return c.Status(400).JSON(
		fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
}

func GetSecrets(logger *zap.Logger) {
	// Load the configuration file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Panic("Error reading config file", zap.Error(err))
		os.Exit(1)
	}

	// Set configuration variables based on struct fields
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Panic("error retrieving secret value", zap.Error(err))
		os.Exit(1)
	}

	viper.Set("DATABASE_URL", config.DatabaseUrl)
	viper.Set("DB_NAME", config.DatabaseName)
	viper.Set("REDIS_URL", config.RedisUrl)
	viper.Set("JWT_SECRET_KEY", config.JWTSecretKey)

}

func GenerateToken(user User) (string, error) {
	usr := map[string]interface{}{
		"email": user.Email,
		"_id":   user.ID.Hex(),
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = usr
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET_KEY")))

	if err != nil {
		return "", err
	}

	return t, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func HandleValidationError(c *fiber.Ctx, err error) error {

	if _, ok := err.(*validator.InvalidValidationError); ok {

		return HandleError(c, err)
	}

	var errMessage string
	for _, err := range err.(validator.ValidationErrors) {
		errMessage = fmt.Sprintf("enter a valid %v in %v field", err.Kind().String(), err.Field())
		break
	}

	return HandleError(c, errors.New(errMessage))
}
