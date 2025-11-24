package itemsvc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/dooleyonline/backend/internal/config"
	_ "golang.org/x/image/webp"
)

func generatePlaceholder(cfg *config.Config, img string) (string, error) {
	publicURL := fmt.Sprintf("%s/%s/%s", cfg.StorageUrl, "item", img)

	res, err := http.Get(publicURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: status %d", res.StatusCode)
	}

	src, err := imaging.Decode(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	thumb := imaging.Thumbnail(src, 5, 5, imaging.Box)

	var buf bytes.Buffer
	if err = imaging.Encode(&buf, thumb, imaging.PNG); err != nil {
		return "", fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", b64), nil
}
