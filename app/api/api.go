package api

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chilts/sid"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type ApiBase struct {
	RedisCli *redis.Client
}

func (e *ApiBase) GetListApiCb(c *gin.Context) {
	keys, err := e.RedisCli.Keys("cb_*").Result()
	if err != nil {
		c.JSON(404, gin.H{
			"result": "File not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"result": keys,
	})
}

func (e *ApiBase) GetListApiCbit(c *gin.Context) {
	keys, err := e.RedisCli.Keys("cbit_*").Result()
	if err != nil {
		c.JSON(404, gin.H{
			"result": "File not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"result": keys,
	})
}

func (e *ApiBase) PostApiCb(c *gin.Context) {
	key := "cb_" + sid.IdBase64()
	err := e.RedisCli.SetNX(key, "this is cb addr", 0).Err()

	if err != nil {
		c.JSON(400, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"result": key,
	})
}

func (e *ApiBase) DeleteApiCb(c *gin.Context) {
	id := c.Param("id")

	val, err := e.RedisCli.Exists(id).Result()
	if err != nil || val == 0 {
		c.JSON(404, gin.H{
			"result": false,
			"error":  "File not found",
		})
		return
	}
	e.RedisCli.Del(id)
	c.JSON(200, gin.H{
		"result": "Ok",
	})
}

func (e *ApiBase) GetOneApiCb(c *gin.Context) {
	id := c.Param("id")
	keys, err := e.RedisCli.Keys("*_" + id).Result()
	if err != nil {
		c.JSON(404, gin.H{
			"result": "File not found",
		})
		return
	}

	var ret []map[string]interface{}
	for _, key := range keys {
		var data map[string]interface{}
		val, err := e.RedisCli.Get(key).Result()
		if err != nil {
			c.JSON(400, gin.H{"result": false, "error": err.Error()})
			return
		}

		if err := json.NewDecoder(strings.NewReader(val)).Decode(&data); err == nil {
			ret = append(ret, data)
		}
	}

	c.JSON(200, gin.H{
		"result": ret,
	})
}

func (e *ApiBase) PostOneApiCb(c *gin.Context) {
	expireMinutes, _ := strconv.ParseInt(os.Getenv("KEY_EXPIRE_MINUTES"), 10, 64)
	id := c.Param("id")
	val, err := e.RedisCli.Exists(id).Result()

	if err != nil || val == 0 {
		c.JSON(404, gin.H{
			"result": "File not found",
		})
		return
	}

	contentType := c.Request.Header.Get("Content-Type")
	var data string
	if contentType == "application/json" {
		var req map[string]interface{}
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			c.Error(err)
			c.JSON(400, gin.H{"result": false, "error": err.Error()})
			return
		}
		tmp, _ := json.Marshal(req)
		data = string(tmp)
	} else {
		c.JSON(400, gin.H{"result": false, "error": "unsupported content type: " + contentType})
	}

	detailKey := "cbit_" + sid.IdBase64() + "_" + id
	err = e.RedisCli.SetNX(detailKey, data, time.Minute*time.Duration(expireMinutes)).Err()
	if err != nil {
		c.JSON(400, gin.H{"result": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"result": "ok",
	})
}
