package stores

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

type DiskImageStore struct {
	// There will be concurrent requests to write
	// to file, so a mutex is needed.
	m           sync.RWMutex          // multiple readers, one writer
	imageFolder string                // folder path
	images      map[string]*ImageInfo // imageId => imageInfo
}

type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	images := make(map[string]*ImageInfo)
	return &DiskImageStore{imageFolder: imageFolder, images: images}
}

func (st *DiskImageStore) Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error) {
	// Generating the image ID.
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Cannot generate image id: %w", err)
	}

	// Generating the image path.
	imagePath := fmt.Sprintf("%s/%s%s", st.imageFolder, imageId, imageType)
	log.Printf("Saving image %s to file %s...", imageId.String(), imagePath)

	// Creating the file at the generated path.
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Cannot create image file: %w", err)
	}

	// Writing the image data to the new file.
	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Cannot write image data to file: %w", err)
	}
	log.Printf("Image %s save to file %s...", imageId.String(), imagePath)

	// Acquiring the lock to update the in-memory
	// counterpart of the image data.
	st.m.Lock() // write lock
	defer st.m.Unlock()

	st.images[imageId.String()] = &ImageInfo{
		LaptopId: laptopId,
		Type:     imageType,
		Path:     imagePath,
	}

	return imageId.String(), nil
}
