package mongodb

import (
	"github.com/Kamva/mgm/v2"
	"github.com/spf13/viper"
	"movies-review-api/domain"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoRepository struct {
	UserRepo domain.UserRepository
}

func New(l *zap.Logger) *MongoRepository {
	// connect to mongodb
	var dbName, connectionString string

	if viper.Get("DB_NAME") != nil {
		dbName = viper.Get("DB_NAME").(string)
	}

	if viper.Get("DATABASE_URL") != nil {
		connectionString = viper.Get("DATABASE_URL").(string)
	}

	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(connectionString), options.Client().SetMaxPoolSize(500))

	if err != nil {
		l.Error(err.Error(), zap.Error(err))
	}

	return &MongoRepository{
		UserRepo: NewUserRepository(l),
	}
}
