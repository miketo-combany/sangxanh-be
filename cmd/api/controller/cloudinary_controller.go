package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type imageController struct {
	imageSvc service.ImageService
}

func NewImageController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &imageController{
		imageSvc: do.MustInvoke[service.ImageService](di),
	}, nil
}

func (c *imageController) Register(g *echo.Group) {
	g = g.Group("/image")
	g.POST("/upload", c.Upload) // POST /image/upload
}

// Upload godoc
// @Summary Upload images
// @Description Upload one or multiple images to Cloudinary
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Image files to upload (multiple files supported)"
// @Param folder formData string false "Cloudinary folder path"
// @Success 200 {object} map[string]interface{} "Upload successful with image URLs"
// @Failure 400 {object} map[string]interface{} "Invalid request or no file uploaded"
// @Router /image/upload [post]
func (c *imageController) Upload(e echo.Context) error {
	// Parse the whole multipart form once. Echo re‑uses http.Request.ParseMultipartForm
	// so this is cached and cheap even for many fields.
	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid multipart form: " + err.Error(),
		})
	}

	// Preferred: <input type="file" name="files" multiple>
	files := form.File["files"]

	// Fallback: single <input type="file" name="file">
	if len(files) == 0 {
		if fh, err := e.FormFile("file"); err == nil {
			files = []*multipart.FileHeader{fh}
		}
	}

	if len(files) == 0 {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": "no file uploaded",
		})
	}

	folder := e.FormValue("folder") // may be ""

	// Pass the slice straight through to the service
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.imageSvc.UploadImages(ctx, files, folder)
	})
}
