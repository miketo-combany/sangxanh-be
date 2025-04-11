package service

import (
	"SangXanh/pkg/common/api"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type ImageService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (api.Response, error)
}

type imageService struct {
	cld *cloudinary.Cloudinary
}

func NewImageService(di do.Injector) (ImageService, error) {
	cld, err := do.Invoke[*cloudinary.Cloudinary](di)
	if err != nil {
		return nil, fmt.Errorf("init ImageService: %w", err)
	}
	return &imageService{cld: cld}, nil
}

func (s *imageService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (api.Response, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	upParams := uploader.UploadParams{
		Folder:   folder,           // "" = root
		PublicID: uuid.NewString(), // random name avoids collisions
	}

	res, err := s.cld.Upload.Upload(ctx, src, upParams) // Cloudinary call
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload: %w", err)
	}

	// res.SecureURL is the HTTPS URL you usually want
	return api.Success(map[string]any{
		"url":       res.SecureURL,
		"public_id": res.PublicID,
		"width":     res.Width,
		"height":    res.Height,
	}), nil
}
