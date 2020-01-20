package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	godaemon "github.com/sevlyar/go-daemon"
)

const ginLogFormat = "request => %d | %s | %s | %s | %s | %s"

type CommonResp struct {
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

type runnerCmd struct {
	Command string `form:"command"`
}

func server(daemon bool) {

	if daemon {
		cntxt := &godaemon.Context{
			PidFileName: "httpcmd.pid",
			PidFilePerm: 0644,
			LogFileName: "httpcmd.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := cntxt.Reborn()
		if err != nil {
			logrus.Fatal("unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer func() { _ = cntxt.Release() }()
		logrus.Info("httpcmd run as daemon mode")
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(ginLog)
	engine.Use(gin.Recovery())
	engine.POST("/", runCmd)

	if err := engine.Run(); err != nil {
		logrus.Fatal(err)
	}
}

func runCmd(c *gin.Context) {

	if token != "" {
		bearerToken := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(bearerToken, "Bearer ") || len(strings.Fields(bearerToken)) != 2 {
			c.JSON(403, gin.H{
				"message":   "invalid token",
				"timestamp": time.Now().Unix(),
			})
			return
		}
		reqToken := strings.Fields(bearerToken)[1]
		if reqToken != token {
			c.JSON(403, gin.H{
				"message":   "invalid token",
				"timestamp": time.Now().Unix(),
			})
			return
		}
	}

	var rcmd runnerCmd

	err := c.ShouldBind(&rcmd)
	if err != nil {
		c.JSON(400, failed("parse parameters failed"))
		return
	}
	logrus.Debugf("request commands: %s", rcmd.Command)

	cmds := strings.Fields(rcmd.Command)
	if len(cmds) == 0 {
		c.JSON(400, failed("command is empty"))
		return
	}

	reg := regexp.MustCompile(cmdRegex)
	if !reg.MatchString(rcmd.Command) {
		c.JSON(400, failed("command not allow"))
		return
	}

	sout, serr, err := run(cmds[0], cmds[1:]...)
	if err != nil {
		logrus.Errorf("%v: %s", err, serr)
		errMsg := fmt.Sprintf("run command failed: %v: %s", err, serr)
		if sout != "" {
			errMsg = errMsg + " output: " + sout
		}
		c.JSON(500, failed(errMsg))
		return
	} else {
		dataMessage := sout
		if serr != "" {
			logrus.Warn(serr)
			dataMessage += "\n" + serr
		}
		c.JSON(200, data(dataMessage))
		return
	}

}

// for the fast return failed result
func failed(message string, args ...interface{}) CommonResp {
	return CommonResp{
		Message:   fmt.Sprintf(message, args...),
		Timestamp: time.Now().Unix(),
	}
}

// for the fast return result with custom data
func data(data interface{}) CommonResp {
	return CommonResp{
		Message:   "success",
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}

func ginLog(c *gin.Context) {
	path := c.Request.URL.Path
	start := time.Now()
	c.Next()
	end := time.Now()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		logrus.Error(fmt.Sprintf(ginLogFormat, c.Writer.Status(), c.ClientIP(), c.Request.Method, path, latency, c.Request.UserAgent()), " | ERROR: ", c.Errors.String())
	} else {
		logrus.Info(fmt.Sprintf(ginLogFormat, c.Writer.Status(), c.ClientIP(), c.Request.Method, path, latency, c.Request.UserAgent()))
	}
}
