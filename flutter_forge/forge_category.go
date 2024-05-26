package flutterforge

import (
	"encoding/json"
	"fmt"
	"forger/db"
	"net/http"

	"github.com/huandu/go-sqlbuilder"
)

type ForgeCategoryModel struct {
	ID          int                     `json:"id"`
	Title       string                  `json:"title"`
	Description *string                 `json:"description"`
	Logo        string                  `json:"logo"`
	Slug        string                  `json:"slug"`
	SubCategory []ForgeSubcategoryModel `json:"forge_subcategory"`
}

type ForgeSubcategoryModel struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	Image           string  `json:"image"`
	ForgeCategoryID int     `json:"forge_category_id"`
	Description     *string `json:"description"`
}

func ForgeCategory(w http.ResponseWriter, r *http.Request) {
	sql := sqlbuilder.Select("id", "title", "description", "logo", "slug").From("flutter_forge.forge_category").String()

	rows, err := db.Database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var categories []*ForgeCategoryModel

	for rows.Next() {
		var id int
		var title, slug string
		var description, logo *string

		if err := rows.Scan(&id, &title, &description, &logo, &slug); err != nil {
			panic(err)
		}

		var logoURL string = fmt.Sprintf("https://flutterforge.s3.ap-south-1.amazonaws.com/forge_category/%s.svg", slug)

		category := &ForgeCategoryModel{
			ID:          id,
			Title:       title,
			Description: description,
			Logo:        logoURL,
			Slug:        slug,
			SubCategory: make([]ForgeSubcategoryModel, 0),
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	sqlSubcategory := sqlbuilder.Select("id", "title", "slug", "image", "forge_category_id", "description").From("flutter_forge.forge_subcategory").String()
	rowsSub, err := db.Database.Query(sqlSubcategory)
	if err != nil {
		panic(err)
	}
	defer rowsSub.Close()

	for rowsSub.Next() {
		var id, forgeCategoryID int
		var title, slug string
		var image, description *string

		if err := rowsSub.Scan(&id, &title, &slug, &image, &forgeCategoryID, &description); err != nil {
			panic(err)
		}

		var imageURL string = fmt.Sprintf("https://flutterforge.s3.ap-south-1.amazonaws.com/forge_subcategory/%s.png", slug)

		if id == 10 || id == 11 || id == 12 {
			imageURL = fmt.Sprintf("https://flutterforge.s3.ap-south-1.amazonaws.com/forge_subcategory/%s.gif", slug)
		}

		for _, category := range categories {
			if category.ID == forgeCategoryID {
				subCategory := ForgeSubcategoryModel{
					ID:              id,
					Title:           title,
					Slug:            slug,
					Image:           imageURL,
					ForgeCategoryID: forgeCategoryID,
					Description:     description,
				}
				category.SubCategory = append(category.SubCategory, subCategory)
			}
		}
	}

	forgeCategoriesJSON, err := json.Marshal(categories)
	if err != nil {
		panic(err)
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(forgeCategoriesJSON)
	if err != nil {
		panic(err)
	}
}
