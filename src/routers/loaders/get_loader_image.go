package loaders

import (
	"context"
	"errors"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// @Tags Погрузчики
// @Summary Получить фото погрузчика по id
// @Description Отправляет фото погрузчика с указанным id
// @ID loader_image
// @Accept json
// @Produce image/jpg, image/png
// @Param loader_id path string true "ID погрузчика" example(1)
// @Success 201 {file} binary
// @Router /api/loaders/getimage/{loader_id} [get]
func GetLoaderImage(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path
	var ctx = context.Background()
	db := connection.GetDatabase()

	loaderID := chi.URLParam(r, "loader_id")

	var loader models.Loader
	err := db.NewSelect().Model(&loader).Where("id = ?", loaderID).Scan(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"invalid data"}`, http.StatusBadRequest)
		return
	}
	imagePath := loader.PicturePath
	_, err = os.Stat(imagePath)
	if errors.Is(err, os.ErrNotExist) {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"picture not exist"}`, http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, imagePath)
	additional.PrintSuccess(path, "Фото погрузчика с ID "+loaderID+" получено")
}
