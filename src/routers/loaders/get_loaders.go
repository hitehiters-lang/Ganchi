package loaders

import (
	"context"
	"encoding/json"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"
)

// @Tags Погрузчики
// @Summary Получить все погрузчики
// @Description Отправляет всю таблицу с данными о каждом погрузчике
// @ID all_loaders
// @Accept json
// @Produce json
// @Success 200 {array} models.Loader
// @Router /api/loaders/getall [get]
func GetLoaders(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path
	var ctx = context.Background()
	var db = connection.GetDatabase()

	var loaders []models.Loader
	err := db.NewSelect().Model(&loaders).Scan(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "loaders not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loaders)

	additional.PrintSuccess(path, "Данные о погрузчиках получены")
}
