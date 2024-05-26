package flutterforge

import (
	"encoding/json"
	"forger/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/huandu/go-sqlbuilder"
)

type ForgeComponentsModel struct {
	AppRoute string `json:"app_route"`
	Gist     string `json:"gist"`
	Title    string `json:"title"`
}

func ForgeComponents(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("forge_components.title", "forge_components.app_route", "forge_components.gist")
	sb.From("flutter_forge.forge_components")
	sb.JoinWithOption(sqlbuilder.InnerJoin, "flutter_forge.forge_subcategory", "forge_components.forge_subcategory_id = forge_subcategory.id")
	sb.Where(sb.Equal("forge_subcategory.slug", slug))
	sqlQuery, args := sb.Build()

	rows, err := db.Database.Query(sqlQuery, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var forgeComponents []ForgeComponentsModel

	for rows.Next() {
		var title, appRoute, gist string
		var components ForgeComponentsModel
		if err := rows.Scan(&title, &appRoute, &gist); err != nil {
			panic(err)
		}

		components = ForgeComponentsModel{
			Title:    title,
			AppRoute: appRoute,
			Gist:     gist,
		}

		forgeComponents = append(forgeComponents, components)

	}

	result, err := json.Marshal(forgeComponents)

	if err != nil {
		panic(err)
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(result)
	if err != nil {
		panic(err)
	}

}
