package repository

import (
	"context"
	"eth_mertics/internal/metrics"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Функция для массового сохранения данных о volume, fees и revenue
func SaveDataBatchVFR(db *pgxpool.Pool, data metrics.DexOverviewVFR, metric string) error {
	// Создаём транзакцию для более эффективного выполнения
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // откатываем транзакцию в случае ошибки

	// Формируем SQL запрос для массовой вставки
	sql := fmt.Sprintf("INSERT INTO %s (name, category, %s) VALUES", metric, metric)
	values := []interface{}{}
	placeholders := []string{}
	for i, v := range data.Protocols {
		// Строим placeholders для вставки
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		values = append(values, v.Name, v.Category, v.Total24h)
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

func SaveDataBatchTvl(db *pgxpool.Pool, data []metrics.ProtocolTvl, metric string) error {
	// Создаём транзакцию для более эффективного выполнения
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // откатываем транзакцию в случае ошибки

	// Формируем SQL запрос для массовой вставки
	sql := fmt.Sprintf("INSERT INTO %s (name, category, %s) VALUES", metric, metric)
	values := []interface{}{}
	placeholders := []string{}
	for i, v := range data {
		// Строим placeholders для вставки
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		values = append(values, v.Name, v.Category, v.Tvl)
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
