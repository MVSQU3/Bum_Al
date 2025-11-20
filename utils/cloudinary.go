package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func initCloudinary() error {
	if cld != nil {
		return nil // Déjà initialisé
	}
	
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return fmt.Errorf("CLOUDINARY_URL n'est pas définie")
	}
	
	var err error
	cld, err = cloudinary.NewFromURL(cloudinaryURL)
	return err
}


func UploadImage(ctx context.Context, filePath string) (string, error) {
	if err := initCloudinary(); err != nil {
		return "", err
	}
	
	resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func UploadImageWithOptions(ctx context.Context, filePath string, publicID string, folder string) (map[string]interface{}, error) {
	if err := initCloudinary(); err != nil {
		return nil, err
	}
	
	overwrite := true
	params := uploader.UploadParams{
		PublicID:  publicID,
		Folder:    folder,
		Overwrite: &overwrite,
	}

	resp, err := cld.Upload.Upload(ctx, filePath, params)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"url":      resp.SecureURL,
		"publicID": resp.PublicID,
		"width":    resp.Width,
		"height":   resp.Height,
	}, nil
}