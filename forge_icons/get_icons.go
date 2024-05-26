package forgeicons

import (
	"encoding/json"
	"fmt"
	"forger/db"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
)

type RequestBody struct {
	IconName       string   `json:"icon_name"`
	IconFamilies   []string `json:"icon_families"`
	IconCategories []string `json:"icon_categories"`
	IconTypes      []string `json:"icon_types"`
}

type Icon struct {
	ID           int    `json:"id"`
	IconFamily   string `json:"icon_family"`
	IconName     string `json:"icon_name"`
	IconType     string `json:"icon_type"`
	IconCategory string `json:"icon_category"`
	CreatedAt    string `json:"created_at"`
	IconPath     string `json:"icon_path"`
}

func GetForgeIcons(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error parsing post data", http.StatusBadRequest)
		return
	}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*").From("forge_icons").Limit(60)

	if reqBody.IconName != "" {
		sb.Where(
			sb.Like("icon_name", "%"+reqBody.IconName+"%"),
		)
	}

	buildQueryFilter(sb, "icon_family", reqBody.IconFamilies)
	buildQueryFilter(sb, "icon_type", reqBody.IconTypes)
	buildQueryFilter(sb, "icon_categorgy", reqBody.IconCategories)

	sql, args := sb.Build()
	rows, err := db.Database.Query(sql, args...)
	if err != nil {
		http.Error(w, "Error executing SQL query", http.StatusInternalServerError)
		log.Fatalf("Error executing SQL query: %v", err)
		return
	}
	defer rows.Close()

	var icons []Icon
	const iconBasePath = "https://raw.githubusercontent.com/devmysip/icons/42d9efe496c453b0c6f295a1d5264ccbe4cae445"
	for rows.Next() {
		var icon Icon
		if err := rows.Scan(&icon.ID, &icon.IconFamily, &icon.IconName, &icon.IconType, &icon.IconCategory, &icon.CreatedAt); err != nil {
			panic(err)
		}
		icon.IconPath = fmt.Sprintf("%s/%s/%s/%s/%s", iconBasePath, icon.IconFamily, icon.IconType, icon.IconCategory, icon.IconName)

		icons = append(icons, icon)
	}

	response := make(map[string]interface{})
	response["status"] = 1
	response["result"] = icons
	response["filters"] = buildFilter()
	if len(icons) == 0 {
		response["result"] = []string{}
	}

	iconsJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error converting icons to JSON", http.StatusInternalServerError)
		log.Fatalf("Error converting icons to JSON: %v", err)
		return
	}

	// Set content type and send JSON response
	w.Header().Set("Content-Type", "application/json")

	w.Write(iconsJSON)
}

func buildQueryFilter(sb *sqlbuilder.SelectBuilder, field string, filters []string) {
	if len(filters) <= 0 {
		return
	}

	filter := make([]string, 0)
	for _, v := range filters {
		if v != "" {
			filter = append(filter, v)
		}
	}

	values := make([]interface{}, len(filter))
	for i, v := range filter {
		values[i] = v
	}
	sb.Where(
		sb.In(field, values...),
	)
}
func buildFilter() map[string]interface{} {

	filter := make(map[string]interface{})

	IconFamilies := []string{
		"Ion",
		"Jam",
		"Line Awesome",
		"Social",
		"Tabler Icons",
		"Unicons",
		"basil",
		"cool icons",
		"iconly",
		"vuaesex",
	}

	iconTypes := []string{
		"Regular",
		"Solid",
		"logo",
		"Filled",
		"Monochrome",
		"Original",
		"Outline",
		"bold",
		"broken",
		"bulk",
		"linear",
		"twotone",
	}

	IconCategories := []string{
		"filled",
		"Brands",
		"Communication",
		"Devices",
		"Files",
		"General",
		"Interface",
		"Media",
		"Navigation",
		"Status",
		"Arrow",
		"Calendar",
		"Edit",
		"File",
		"Live",
		"Menu",
		"System",
		"User",
		"Warning",
		"Bold",
		"Broken",
		"Bulk",
		"Curved",
		"Light",
		"Light-Outline",
		"Sharp",
		"Two-tone",
		"Archive",
		"Astrology",
		"Building",
		"Business",
		"Car",
		"Computers-Devices-Electronics",
		"Content-Edit",
		"Crypto-Company",
		"Delivery",
		"Design-Tools",
		"Emails-Messages",
		"Essetional",
		"Grid",
		"Location",
		"Money",
		"Notifications",
		"Programing",
		"School-Learning",
		"Search",
		"Security",
		"Settings",
		"Shop",
		"Support-Like-Question",
		"Time",
		"Type-Paragraph-Character",
		"Users",
		"Video-Audio-Image",
		"Weather",
		"Call",
		"Crypto-Currency",
	}

	filter["icon_families"] = IconFamilies
	filter["icon_types"] = iconTypes
	filter["icon_categories"] = IconCategories

	return filter
}
