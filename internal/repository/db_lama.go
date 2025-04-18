package repository

import (
	"context"
	"eth_mertics/internal/lama"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Функция для массового сохранения данных о volume, fees и revenue
func SaveDataBatchVFR(db *pgxpool.Pool, data []lama.ProtocolInfo, metric string) error {
	if len(data) == 0 {
		return nil
	}

	// Собираем уникальные chains из data
	chainSet := make(map[string]struct{})
	for _, v := range data {
		for chain := range v.Breakdown24 {
			chainSet[chain] = struct{}{}
		}
	}

	// Конвертируем в массив и сортируем для стабильности
	var chains []string
	for chain := range chainSet {
		chains = append(chains, chain)
	}
	sort.Strings(chains)

	// Подготовка к вставке
	columns := append([]string{"name", "category", metric}, chains...)

	const maxParams = 65535
	fieldsPerRow := len(columns)
	maxRows := maxParams / fieldsPerRow

	tx, err := db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Разбивка на батчи
	for start := 0; start < len(data); start += maxRows {
		end := start + maxRows
		if end > len(data) {
			end = len(data)
		}

		var (
			sqlBuilder   strings.Builder
			values       []interface{}
			placeholders []string
		)

		sqlBuilder.WriteString(fmt.Sprintf(`INSERT INTO %s (`, metric))
		sqlBuilder.WriteString(strings.Join(columns, ", "))
		sqlBuilder.WriteString(") VALUES ")

		for _, v := range data[start:end] {
			ph := []string{}
			offset := len(values)

			// name, category, total24h
			values = append(values, v.Name, v.Category, v.Total24h)
			ph = append(ph, fmt.Sprintf("$%d", offset+1), fmt.Sprintf("$%d", offset+2), fmt.Sprintf("$%d", offset+3))

			// По каждому чейну
			for _, chain := range chains {
				val, ok := v.Breakdown24[chain]
				if ok {
					values = append(values, val)
				} else {
					values = append(values, nil)
				}
				ph = append(ph, fmt.Sprintf("$%d", len(values)))
			}

			placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(ph, ", ")))
		}

		sqlBuilder.WriteString(strings.Join(placeholders, ", "))

		if _, err := tx.Exec(context.Background(), sqlBuilder.String(), values...); err != nil {
			return fmt.Errorf("insert exec error: %w", err)
		}
	}

	return tx.Commit(context.Background())
}

func SaveDataBatchTvl(db *pgxpool.Pool, data []lama.TvlChains) error {
	// Создаём транзакцию для более эффективного выполнения
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // откатываем транзакцию в случае ошибки

	// Формируем SQL запрос для массовой вставки
	sql := "INSERT INTO tvl_Chains (name, category, tvl) VALUES"
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

func SaveDataBatchTvlProtocols(db *pgxpool.Pool, data []lama.ProtocolTvl) error {
	if len(data) == 0 {
		return nil
	}

	// Собираем уникальные chains из data
	chainSet := make(map[string]struct{})
	for _, v := range data {
		for chain := range v.ChainTVLs {
			chainSet[chain] = struct{}{}
		}
	}

	// Конвертируем в массив и сортируем для стабильности
	var chains []string
	for chain := range chainSet {
		chains = append(chains, chain)
	}
	sort.Strings(chains)

	// Подготовка к вставке
	columns := append([]string{"name", "category", "tvl"}, chains...)

	const maxParams = 65535
	fieldsPerRow := len(columns)
	maxRows := maxParams / fieldsPerRow

	tx, err := db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Разбивка на батчи
	for start := 0; start < len(data); start += maxRows {
		end := start + maxRows
		if end > len(data) {
			end = len(data)
		}

		var (
			sqlBuilder   strings.Builder
			values       []interface{}
			placeholders []string
		)

		sqlBuilder.WriteString(`INSERT INTO tvl_Protocols (`)
		sqlBuilder.WriteString(strings.Join(columns, ", "))
		sqlBuilder.WriteString(") VALUES ")

		for _, v := range data[start:end] {
			ph := []string{}
			offset := len(values)

			// name, category, total24h
			values = append(values, v.Name, v.Category, v.Tvl)
			ph = append(ph, fmt.Sprintf("$%d", offset+1), fmt.Sprintf("$%d", offset+2), fmt.Sprintf("$%d", offset+3))

			// По каждому чейну
			for _, chain := range chains {
				val, ok := v.ChainTVLs[chain]
				if ok {
					values = append(values, val)
				} else {
					values = append(values, nil)
				}
				ph = append(ph, fmt.Sprintf("$%d", len(values)))
			}

			placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(ph, ", ")))
		}

		sqlBuilder.WriteString(strings.Join(placeholders, ", "))

		if _, err := tx.Exec(context.Background(), sqlBuilder.String(), values...); err != nil {
			return fmt.Errorf("insert exec error: %w", err)
		}
	}

	return tx.Commit(context.Background())
}

func SaveDataBatchMcapProtocols(db *pgxpool.Pool, data []lama.ProtocolTvl) error {
	// Создаём транзакцию для более эффективного выполнения
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // откатываем транзакцию в случае ошибки

	// Формируем SQL запрос для массовой вставки
	sql := "INSERT INTO Mcap_Protocols (name, category, mcap) VALUES"
	values := []interface{}{}
	placeholders := []string{}
	for i, v := range data {
		// Строим placeholders для вставки
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		values = append(values, v.Name, v.Category, v.Mcap)
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
