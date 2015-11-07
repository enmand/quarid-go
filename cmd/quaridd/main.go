package main

import (
	"os"

	"github.com/enmand/quarid-go/pkg/config"
	"github.com/enmand/quarid-go/pkg/http"
	"github.com/enmand/quarid-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

func init() {
	c := config.Get()
	env := c.GetString("env")
	if env == "" {
		env = os.Getenv("GIN_MODE")
	}

	if env != "" {
		gin.SetMode(env)
	}
}

func main() {
	logger.Log.Info("Loading Quarid...")
	c := config.Get()

	r := http.Router(
		gin.Logger(),
		gin.Recovery(),
	)
	r.Run(c.GetString("listen"))
}
