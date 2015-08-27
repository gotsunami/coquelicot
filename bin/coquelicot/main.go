package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gotsunami/coquelicot"
)

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("version: %s\n", appVersion)
		return
	}

	s := coquelicot.NewStorage(*storage)

	r := gin.Default()
	r.Use(coquelicot.CORSMiddleware())
	r.Use(static.ServeRoot("/", s.StorageDir()))

	r.POST("/files", s.UploadHandler)

	log.Printf("Storage place in: %s", s.StorageDir())
	log.Printf("Start server on %s", *host)
	r.Run(*host)

}
