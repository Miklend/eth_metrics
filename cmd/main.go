package main

import (
	"eth_mertics/internal/coingecko"
	"eth_mertics/internal/lama"
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

	lamaTvlProtocols(db)
	lamaTvlChains(db)
	lamaVFR(db)
	coingeckoVolMcap(db)

}

func coingeckoVolMcap(db *pgxpool.Pool) {
	data, err := coingecko.GetDataVolMcapCoingecko(os.Getenv("urlCoingeckoVolMcap"), os.Getenv("ApiKeyCoinGecko"))
	if err != nil {
		logrus.Fatalf("failed to get data: %s", err.Error())
	}
	err = repository.SaveDataBatchCoingecko(db, data)
	if err != nil {
		logrus.Fatalf("failed to save data: %s", err.Error())
	}
	logrus.Infof("Successfully saved coingecko data to the database")
}

func lamaVFR(db *pgxpool.Pool) {
	for metric, url := range map[string]string{
		"volume24":  os.Getenv("urlVolume"),
		"fees24":    os.Getenv("urlFees"),
		"revenue24": os.Getenv("urlRevenue"),
	} {
		data, err := lama.GetDataVFR(url)
		if err != nil {
			logrus.Fatalf("failed to get data: %s", err.Error())
		}
		err = repository.SaveDataBatchVFR(db, data, metric)
		if err != nil {
			logrus.Fatalf("failed to save data: %s", err.Error())
		}
		logrus.Infof("Successfully saved %s data to the database", metric)
	}
}

func lamaTvlChains(db *pgxpool.Pool) {
	data, err := lama.GetDataTvlChains(os.Getenv("urlTvlchains"))
	if err != nil {
		logrus.Fatalf("failed to get data: %s", err.Error())
	}
	err = repository.SaveDataBatchTvl(db, data)
	if err != nil {
		logrus.Fatalf("failed to save data: %s", err.Error())
	}
	logrus.Infof("Successfully saved tvlChains data to the database")
}

func lamaTvlProtocols(db *pgxpool.Pool) {
	dataTvl, dataMcap, err := lama.GetDataTvlProtocols(os.Getenv("urlTvlProtocols"))
	if err != nil {
		logrus.Fatalf("failed to get data: %s", err.Error())
	}
	err = repository.SaveDataBatchTvlProtocols(db, dataTvl)
	if err != nil {
		logrus.Fatalf("failed to save data: %s", err.Error())
	}
	logrus.Infof("Successfully saved tvlProtocols data to the database")
	err = repository.SaveDataBatchMcapProtocols(db, dataMcap)
	if err != nil {
		logrus.Fatalf("failed to save data: %s", err.Error())
	}
	logrus.Infof("Successfully saved tvlMcap data to the database")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
