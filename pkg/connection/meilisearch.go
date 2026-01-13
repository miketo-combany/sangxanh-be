package connection

import (
	"SangXanh/pkg/enum"
	"SangXanh/pkg/log"
	"fmt"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
)

// This small command connects to Supabase, fetches all
// non-deleted products and logs them to the console.
func SyncData(di do.Injector) {
	// Initialize dependency injection and shared components

	// Get MeiliSearch configuration from environment variables
	meilisearchURL := os.Getenv("MEILISEARCH_URL")
	log.Info("MeiliSearch URL: ", meilisearchURL)
	if meilisearchURL == "" {
		meilisearchURL = "http://127.0.0.1:7700" // fallback for local dev
	}
	meilisearchKey := os.Getenv("MEILISEARCH_API_KEY")
	if meilisearchKey == "" {
		meilisearchKey = "masterKey"
	}
	meilisearchClient := meilisearch.New(meilisearchURL, meilisearch.WithAPIKey(meilisearchKey))

	// Get Supabase client from the DI container
	client := do.MustInvoke[*supabase.Client](di)

	type Product struct {
		Id           string              `json:"id"`
		Name         string              `json:"name"`
		Price        float32             `json:"price"`
		Content      string              `json:"content"`
		ImageDetail  string              `json:"image_detail"`
		Thumbnail    string              `json:"thumbnail"`
		CategoryId   string              `json:"category_id"`
		Discount     float32             `json:"discount"`
		DiscountType enum.DiscountType   `json:"discount_type"`
		ProductCode  string              `json:"product_code"`
		Description  string              `json:"description"`
		Metadata     map[string]string   `json:"metadata"`
		CreatedAt    time.Time           `json:"created_at"`
		UpdatedAt    time.Time           `json:"updated_at"`
		DeletedAt    time.Time           `json:"deleted_at"`
		Questions    []map[string]string `json:"questions"`
	}

	var products []Product
	if err := client.DB.
		From("products").
		Select("id,name,price,content,description,category_id,thumbnail,image_detail,discount,discount_type,product_code,metadata,created_at,updated_at,deleted_at,questions").
		IsNull("deleted_at").
		Execute(&products); err != nil {
		log.Fatal("failed to fetch products:", err)
	}

	log.Infof("Fetched %d products", len(products))

	meilisearchClient.Index("products").DeleteAllDocuments(nil)
	log.Info("Delete all documents successfully")
	var meiliDocs []map[string]interface{}
	for _, p := range products {
		doc := map[string]interface{}{
			"id":    p.Id,
			"name":  p.Name,
			"price": p.Price,
			// "content":       p.Content,
			// "description":   p.Description,
			// "category_id":   p.CategoryId,
			// "thumbnail":     p.Thumbnail,
			// "image_detail":  p.ImageDetail,
			// "discount":      p.Discount,
			// "discount_type": p.DiscountType,
			// "product_code":  p.ProductCode,
			// "metadata":      p.Metadata,
			// "questions":     p.Questions,
		}
		meiliDocs = append(meiliDocs, doc)
	}
	id := "id"
	task, err := meilisearchClient.Index("products").AddDocuments(meiliDocs, &meilisearch.DocumentOptions{
		PrimaryKey: &id,
	})

	// updateIndexTask, err := meilisearchClient.Index("products").UpdateFilterableAttributes(&[]interface{}{
	// 	"category_id",
	// 	"discount",
	// 	"price",
	// 	"discount_type",
	// 	"product_code",
	// })

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Init task, ", task.TaskUID)
	// fmt.Println("Update filterable attributes task, ", updateIndexTask.TaskUID)

	log.Info("Sync completed successfully")
}
