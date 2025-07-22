package handlers

import (
	"encoding/json"
	"net/http"
	database "sub/internal/db"
	"sub/internal/models"
	"sub/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// CreateUser добавляет нового пользователя
// @Summary Добавление пользователя
// @Description Добавляет нового пользователя с подпиской
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserSub true "Информация о подписке"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Ошибка валидации или получения данных"
// @Failure 500 {string} string "Ошибка при добавлении пользователя"
// @Router /api/user [post]
func CreateUser(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us, err := ReadBody(r)
		if err != nil {
			http.Error(w, "Ошибка при получении данных: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = database.AddUserToDb(db, us)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// CreateUser Поиск пользователя по user_id
// @Summary Поиск пользователя
// @Description Поиск пользователя по user_id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserSub
// @Failure 400 {string} string "Ошибка валидации или получения данных"
// @Failure 500 {string} string "Ошибка при получении пользователя"
// @Router /api/user/{id} [get]
func GetUser(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			http.Error(w, "Не передан user_id", http.StatusBadRequest)
			return
		}

		var userSub models.UserSub
		err := db.Get(&userSub, `SELECT * FROM user_subscriptions WHERE user_id = $1 ORDER BY start_date DESC LIMIT 1`, userID)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userSub)
	}
}

// UpdateUser обновляет информацию о пользователе
// @Summary Обновление информации о пользователе
// @Description Обновляет информацию о подписке пользователя по user_id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UserSub true "Информация о подписке"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Ошибка валидации или получения данных"
// @Failure 500 {string} string "Ошибка при обновлении пользователя"
// @Router /api/user/{id} [put]
func UpdateUser(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			http.Error(w, "Не передан user_id", http.StatusBadRequest)
			return
		}

		var data models.UserSub

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Ошибка при декодировании данных: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := database.UpdateUserInDb(db, userID, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// List Получение информации всех пользователей
// @Summary Получение информации всех пользователей
// @Description Получение информации всех пользователей
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.UserSub
// @Failure 400 {string} string "Ошибка валидации или получения данных"
// @Router /api/user [get]
func GetAllUsers(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.UserSub
		err := db.Select(&users, `SELECT * FROM user_subscriptions ORDER BY id ASC`)
		if err != nil {
			http.Error(w, "Ошибка при получении пользователей: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// DeleteUser удаляет пользователя по user_id
// @Summary Удаление пользователя
// @Description Удаляет пользователя по user_id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "Пользователь успешно удален"
// @Failure 400 {string} string "Не передан user_id"
// @Failure 500 {string} string "Ошибка при удалении пользователя"
// @Router /api/user/{id} [delete]
func DeleteUser(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			http.Error(w, "Не передан user_id", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`DELETE FROM user_subscriptions WHERE user_id = $1`, userID)
		if err != nil {
			http.Error(w, "Ошибка при удалении пользователя: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetTotalPrice(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")
		userID := r.URL.Query().Get("user_id")
		serviceName := r.URL.Query().Get("service_name")

		var startDate models.MonthYear
		var endDate *models.MonthYear

		if err := startDate.UnmarshalJSON([]byte(startDateStr)); err != nil {
			http.Error(w, "Некорректный формат start_date: "+err.Error(), http.StatusBadRequest)
			return
		}

		if endDateStr != "" {
			var tmp models.MonthYear
			if err := tmp.UnmarshalJSON([]byte(endDateStr)); err != nil {
				http.Error(w, "Некорректный формат end_date: "+err.Error(), http.StatusBadRequest)
				return
			}
			endDate = &tmp
		}

		us := models.UserSub{
			StartDate: startDate,
			EndDate:   endDate,
			UserId:    userID,
			ServiceName: serviceName,
		}

		totalPrice, err := database.CalculateTotalPrice(db, us)
		if err != nil {
			http.Error(w, "Ошибка при расчете общей стоимости: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]int{"total_price": totalPrice}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Ошибка при формировании ответа: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ReadBody(r *http.Request) (models.UserSub, error) {
	var us models.UserSub
	err := json.NewDecoder(r.Body).Decode(&us)
	if err != nil {
		return models.UserSub{}, err
	}

	if utils.IsStructEmpty(us) {
		return models.UserSub{}, err
	}

	if err = utils.Valid(us); err != nil {
		return models.UserSub{}, err
	}
	return us, err
}