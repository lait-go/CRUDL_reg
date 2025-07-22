package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sub/internal/models"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func ConnectDb(key string) (*sqlx.DB, error) {
	if key == "" {
		return nil, errors.New("Ошибка при подключении к базе данных, проверьте конфиги")
	}

	runMigrations(key)

	db, err := sqlx.Connect("postgres", key)
	if err != nil {
		return nil, errors.New("Ошибка при подключении к базе данных")
	}

	return db, nil
}

func runMigrations(key string) {
	m, err := migrate.New("file://migrations", key)

	if err != nil {
		log.Fatal("Ошибка инициализации миграций:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Ошибка выполнения миграций:", err)
	}
}

func AddUserToDb(db *sqlx.DB, us models.UserSub) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var end_date time.Time
	if us.EndDate != nil {
		end_date = us.EndDate.ToTime()
	} else {
		end_date = time.Time{}
	}

	_, err := db.ExecContext(ctx, `INSERT INTO user_subscriptions (
				service_name,
				monthly_price,
				user_id,
				start_date,
				end_date
			) VALUES ($1, $2, $3, $4, $5);`,
		us.ServiceName, us.MonthlyPrice, us.UserId, us.StartDate.ToTime(), end_date)
	if err != nil {
		return errors.New("Ошибка при добавлении подписки:" + err.Error())
	}

	return nil
}

func UpdateUserInDb(db *sqlx.DB, user_id string, us models.UserSub) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var userData models.UserSub

	if user_id == "" {
		return errors.New("ID пользователя не указан")
	}else if us.UserId != "" && us.UserId != user_id {
		return errors.New("ID пользователя не подлежит изменению")
	}

	err := db.QueryRowContext(ctx, `SELECT * FROM user_subscriptions WHERE user_id = $1;`,
		user_id).Scan(&userData.Id, &userData.ServiceName, &userData.MonthlyPrice, &userData.UserId, &userData.StartDate, &userData.EndDate)
	if err != nil {
		return errors.New("Ошибка при получении данных пользователя: " + err.Error())
	}

	if us.ServiceName != "" {
		userData.ServiceName = us.ServiceName
	}
	if us.MonthlyPrice > 0 {
		userData.MonthlyPrice = us.MonthlyPrice
	}
	if us.StartDate != (models.MonthYear{}) {
		userData.StartDate = us.StartDate
	}

	var end_date time.Time
	if us.EndDate != nil {
		end_date = us.EndDate.ToTime()
	} else {
		end_date = userData.EndDate.ToTime()
	}

	_, err = db.ExecContext(ctx, `UPDATE user_subscriptions SET
				service_name = $1,
				monthly_price = $2,
				start_date = $3,
				end_date = $4
			WHERE user_id = $5;`,
		userData.ServiceName, userData.MonthlyPrice, userData.StartDate.ToTime(), end_date, user_id)

	if err != nil {
		return errors.New("Ошибка при обновлении подписки:" + err.Error())
	}

	return nil
}

func CalculateTotalPrice(db *sqlx.DB, us models.UserSub) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var totalPrice int

	conditions := []string{
		"start_date <= $1",
		"(end_date = '0001-01-01' OR end_date >= $2)",
	}

	args := []interface{}{us.EndDate.ToTime(), us.StartDate.ToTime()}

	argID := 3

	if us.UserId != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argID))
		args = append(args, us.UserId)
		argID++
	}

	if us.ServiceName != "" {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argID))
		args = append(args, us.ServiceName)
		argID++
	}

	query := fmt.Sprintf(`
		SELECT SUM(monthly_price)
		FROM user_subscriptions
		WHERE %s;
	`, strings.Join(conditions, " AND "))

	err := db.QueryRowContext(ctx, query, args...).Scan(&totalPrice)
	if err != nil {
		return 0, fmt.Errorf("ошибка при расчете общей стоимости: %w", err)
	}

	return totalPrice, nil
}
