package coquelicot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

// UploadHandler is the endpoint for uploading and storing files.
func (s *Storage) UploadHandler(c *gin.Context) {
	converts, err := GetConvertParams(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("Query params: %s", err),
		})
		return
	}
	converts["original"] = ""

	pavo, _ := c.Request.Cookie("pavo")
	if pavo == nil {

		pavo = &http.Cookie{
			Name:    "pavo",
			Value:   uuid.New(),
			Expires: time.Now().Add(10 * 356 * 24 * time.Hour),
			Path:    "/",
		}
		c.Request.AddCookie(pavo)
		http.SetCookie(c.Writer, pavo)
	}

	files, err := Process(c.Request, s.StorageDir())
	if err == Incomplete {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
			"file":   gin.H{"size": files[0].Size},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusText(http.StatusBadRequest),
			"error":  fmt.Sprintf("Upload error: %q", err.Error()),
		})
		return
	}

	data := make([]map[string]interface{}, 0)
	// Expected status if no error
	status := http.StatusCreated

	for _, ofile := range files {
		attachment, err := Create(s.StorageDir(), ofile, converts)
		if err != nil {
			data = append(data, map[string]interface{}{
				"name":  ofile.Filename,
				"size":  ofile.Size,
				"error": err.Error(),
			})
			status = http.StatusInternalServerError
			continue
		}
		data = append(data, attachment.ToJson())
	}

	c.JSON(status, gin.H{"status": http.StatusText(status), "files": data})
}

// Get parameters for convert from Request query string
func GetConvertParams(req *http.Request) (map[string]string, error) {
	raw_converts := req.URL.Query().Get("converts")

	if raw_converts == "" {
		raw_converts = "{}"
	}

	convert := make(map[string]string)

	err := json.Unmarshal([]byte(raw_converts), &convert)
	if err != nil {
		return nil, err
	}

	return convert, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		// c.Next()
	}
}
