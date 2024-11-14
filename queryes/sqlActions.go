package queryes

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"telegram/queryes/sql"
	"time"
)

func UpdateData(conn *pgx.Conn) error {
	// Обновляем поле date в таблице test
	currentTime := time.Now()
	formattedDate := currentTime.Format("060102")
	sqlFilePath := "./queryes/sql/update_date.sql"
	_, err := conn.Exec(context.Background(), sql.ReadSQLFile(sqlFilePath), formattedDate)
	return err
}

func GetData(ctx context.Context, conn *pgx.Conn) ([]map[string]interface{}, error) {
	select {
	case <-time.After(time.Second * 2):
		filePath := "./queryes/sql/select_date.sql"
		rows, err := conn.Query(context.Background(), sql.ReadSQLFile(filePath))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// Get columns names
		columns := rows.FieldDescriptions()

		// Create slice for results
		var results []map[string]interface{}

		for rows.Next() {
			// Make slice for save data of each row
			values := make([]interface{}, len(columns))
			for i := range values {
				values[i] = new(interface{})
			}

			// read the values into a slice
			if err := rows.Scan(values...); err != nil {
				return nil, err
			}

			// Make map for save data with name columns
			row := make(map[string]interface{})
			for i, field := range columns {
				row[string(field.Name)] = *(values[i].(*interface{}))
			}

			// Add row to result
			results = append(results, row)
		}

		return results, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func TestContext(ctx context.Context) {
	//time.Sleep(time.Second * 3)
	select {
	case <-time.After(7 * time.Second):
		// Если операция завершилась успешно
		fmt.Println("Context test completed")
	case <-ctx.Done():
		// Если контекст был отменен
		fmt.Println("Context test canceled:", ctx.Err())
	}
}
