package loaders

import (
	"context"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// @Tags Погрузчики
// @Summary Получить все погрузчики
// @Description Отправляет всю таблицу с данными о каждом погрузчике
// @ID loader_image
// @Accept json
// @Produce image/jpg, image/png
// @Param loader_id path string true "ID погрузчика" example(1)
// @Success 201 {file} binary
// @Router /api/loaders/getimage/{loader_id} [get]
func GetLoaderImage(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()
	db := connection.GetDatabase()

	router := "/api/loaders/getall"
	loaderID := chi.URLParam(r, "loader_id")

	var loader models.Loader
	db.NewSelect().Model(&loader).Where("id = ?", loaderID).Scan(ctx)

	image := loader.PicturePath

	http.ServeFile(w, r, image)
	additional.PrintSuccess(router, "Фото погрузчика с ID "+loaderID+" получено")
}
