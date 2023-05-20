package handlers_api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func GetWidgetByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := api.DB.GetWidget(widgetID)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}

	out, err := json.MarshalIndent(widget, "", "	")
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(out)
}
