package repository

import (
	"context"
	"eth_mertics/internal/coingecko"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveDataBatchCoingecko(db *pgxpool.Pool, data []coingecko.CurrencyData) error {
	// Создаём транзакцию для более эффективного выполнения
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // откатываем транзакцию в случае ошибки

	// Формируем SQL запрос для массовой вставки
	sql := "INSERT INTO coingecko (created_at, name, category, volume24, mcap) VALUES"
	values := []interface{}{}
	placeholders := []string{}
	for i, v := range data {
		// Строим placeholders для вставки
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", 5*i+1, 5*i+2, 5*i+3, 5*i+4, 5*i+5))
		values = append(values, v.LastUpdatedAt, v.Name, "Chain", v.USD24hVol, v.USDMarketCap)
	}

	// Добавляем placeholders в SQL запрос
	sql += " " + strings.Join(placeholders, ", ")

	// Выполняем запрос
	_, err = tx.Exec(context.Background(), sql, values...)
	if err != nil {
		return err
	}

	// Фиксируем транзакцию
	return tx.Commit(context.Background())
}
