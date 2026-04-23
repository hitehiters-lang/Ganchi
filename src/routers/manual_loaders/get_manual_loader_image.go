package manual_loaders

import (
	"database/sql"
	"errors"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
)

// @Tags Ручные погрузчики
// @Summary Получить фото ручного погрузчика по id
// @Description Отправляет фото погрузчика с указанным id
// @ID get_manual_loader_image
// @Produce image/png
// @Param loader_id path string true "ID ручного погрузчика" example(1)
// @Success 200 {file} binary
// @Router /api/manual_loaders/getimage/{loader_id} [get]
func GetManualLoaderImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()
	db := connection.GetDatabase()

	loaderID := chi.URLParam(r, "loader_id")

	var loader models.ManualLoader
	err := db.NewSelect().Model(&loader).Where("id = ?", loaderID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"loader not found"}`, http.StatusBadRequest)
		return
	}
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	var imagePath string
	switch runtime.GOOS {
	case "windows":
		imagePath = loader.PicturePathWindows
	default:
		imagePath = loader.PicturePathLinux
	}
	_, err = os.Stat(imagePath)
	if errors.Is(err, os.ErrNotExist) {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"picture not exist"}`, http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, imagePath)
	additional.PrintSuccess(path, "Фото погрузчика с ID "+loaderID+" получено")
}
