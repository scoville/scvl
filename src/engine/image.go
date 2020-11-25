package engine

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"mime/multipart"
	"net/http"

	"github.com/scoville/scvl/src/domain"

	"image/gif"
	"image/jpeg"
	"image/png"

	"os"
)

// UploadImageRequest is the request struct for the UploadImage function
type UploadImageRequest struct {
	File     multipart.File
	FileName string
	UserID   int
}

// UploadImage uploads an image
func (e *Engine) UploadImage(req UploadImageRequest) (dimg *domain.Image, err error) {
	img, ext, err := image.Decode(req.File)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	prefix := ""
	switch ext {
	case "png":
		prefix = "data:image/png;base64,"
		err = png.Encode(buf, img)
	case "jpg":
		fallthrough
	case "jpeg":
		prefix = "data:image/jpeg;base64,"
		err = jpeg.Encode(buf, img, nil)
	case "gif":
		prefix = "data:image/gif;base64,"
		err = gif.Encode(buf, img, nil)
	default:
		err = errors.New("unknown image format")
	}
	if err != nil {
		return
	}
	imgBase64Str := prefix + base64.StdEncoding.EncodeToString(buf.Bytes())

	resp, err := http.Post("https://"+os.Getenv("IMAGE_DOMAIN"), "application/json", bytes.NewReader([]byte(`{"image": "`+imgBase64Str+`"}`)))
	if err != nil {
		return
	}
	var res struct {
		URL     string `json:"image_url"`
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return
	}
	if res.URL == "" {
		err = errors.New("failed to upload image: " + res.Message)
		return
	}
	dimg = &domain.Image{
		UserID: req.UserID,
		URL:    res.URL,
	}
	err = e.sqlClient.CreateImage(dimg)
	return
}
