package loaders

import (
	"context"
	"encoding/json"
	"ganchi_app/additional"
	"ganchi_app/config"
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
// @Param lang query string false "Языковой код"
// @Success 200 {array} models.Loader
// @Router /api/loaders/getall [get]
func GetLoaders(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "ru"
	}
	allowedLangs := map[string]bool{"ru": true, "en": true, "kk": true, "zh": true}
	if !allowedLangs[lang] {
		lang = "ru"
	}

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
	var localisations []models.Localisation
	err = db.NewSelect().Model(&localisations).Scan(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "languages not found"}`, http.StatusNotFound)
		return
	}

	locMap := make(map[string]map[string]string, len(localisations))
	for _, loc := range localisations {
		locMap[loc.Loc] = map[string]string{
			"ru": loc.Ru,
			"en": loc.En,
			"kk": loc.Kk,
			"zh": loc.Zh,
		}
	}

	for i := range loaders {
		if translations, ok := locMap[loaders[i].EngineType]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].EngineType = translated
			}
		}
		if translations, ok := locMap[loaders[i].WheelAxis]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].WheelAxis = translated
			}
		}
		if translations, ok := locMap[loaders[i].BrakeType]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].BrakeType = translated
			}
		}
		if translations, ok := locMap[loaders[i].LiftingCylinder]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].LiftingCylinder = translated
			}
		}
		if translations, ok := locMap[loaders[i].HydraulicLiftingEngine]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].HydraulicLiftingEngine = translated
			}
		}
		if translations, ok := locMap[loaders[i].SteeringMode]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].SteeringMode = translated
			}
		}
		if translations, ok := locMap[loaders[i].ChargingTime]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].ChargingTime = translated
			}
		}
		if translations, ok := locMap[loaders[i].WorkingHours]; ok {
			if translated, ok := translations[lang]; ok && translated != "" {
				loaders[i].WorkingHours = translated
			}
		}
		switch lang {
		case "en":
			loaders[i].Price = int(float64(loaders[i].Price) * config.Dol)
		case "kk":
			loaders[i].Price = int(float64(loaders[i].Price) * config.Ten)
		case "zh":
			loaders[i].Price = int(float64(loaders[i].Price) * config.Yua)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loaders)

	additional.PrintSuccess(path, "Данные о погрузчиках на языке "+lang+" получены")
}
