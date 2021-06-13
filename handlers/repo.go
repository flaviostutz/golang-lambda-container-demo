package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *HTTPServer) setupRepoHandlers() {
	h.Router.POST("/repo/:key", createKey(h))
	h.Router.GET("/repo", listKeys(h))
}

func createKey(h *HTTPServer) func(*gin.Context) {
	return func(c *gin.Context) {

		if h.Opt.Readonly {
			c.JSON(http.StatusForbidden, gin.H{"message": "This repo is readonly"})
			return
		}

		key := c.Param("key")

		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "'key' name cannot be empty"})
			return
		}

		logrus.Infof("Creating key %s", c.Param("key"))

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logrus.Infof("createKey: Error reading body contents. err=%s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading body contents"})
			return
		}

		h.Repo[key] = string(body)
	}
}

func listKeys(h *HTTPServer) func(*gin.Context) {
	return func(c *gin.Context) {
		logrus.Infof("Listing all key values")
		c.JSON(200, h.Repo)
	}
}
