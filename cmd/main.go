package main

import (
	"eth_mertics/internal/metrics"
	"eth_mertics/internal/repository"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()
	logrus.Info("Successfully connected to the database")

	lamaVFR(db)
	lamaTvl(db)

}

func lamaVFR(db *pgxpool.Pool) {
	for metric, url := range map[string]string{
		"volume":  os.Getenv("urlVolume"),
		"fees":    os.Getenv("urlFees"),
		"revenue": os.Getenv("urlRevenue"),
	} {
		data, err := metrics.GetDataVFR(url)
		if err != nil {
			logrus.Fatalf("failed to get data: %s", err.Error())
		}
		err = repository.SaveDataBatchVFR(db, *data, metric)
		if err != nil {
			logrus.Fatalf("failed to save data: %s", err.Error())
		}
		logrus.Infof("Successfully saved %s data to the database", metric)
	}
}

func lamaTvl(db *pgxpool.Pool) {
	data, err := metrics.GetDataTvl(os.Getenv("urlTvlProtocols"), os.Getenv("urlTvlchains"))
	if err != nil {
		logrus.Fatalf("failed to get data: %s", err.Error())
	}
	err = repository.SaveDataBatchTvl(db, data, "tvl")
	if err != nil {
		logrus.Fatalf("failed to save data: %s", err.Error())
	}
	logrus.Infof("Successfully saved tvl data to the database")
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
