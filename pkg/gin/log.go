/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 22/01/2021
*/

package gin

import (
	"errors"
	"net/http"
	"strings"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/w6d-io/ci-status/internal/config"
)

const (
	CORRELATION_ID string = "Correlation-ID"
)

func NewHandler() {

	engine := gin.New()
	engine.Use(JSONLogMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(CorrelationID())
	if config.IsAuth() {
		//engine.Use(Auth())
	}
	//engine.POST("/route53/:action",Hdl.Route53)
	//Hdl.Engine = engine
}

// JSONLogMiddleware logs a gin HTTP request in JSON format, with some additional custom key/values
func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := ctrl.Log.WithName("gin")
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := GetDurationInMilliseconds(start)
		entry := log.WithValues(
			"client_ip", GetClientIP(c),
			"duration", duration,
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"referrer", c.Request.Referer(),
			"request_id", c.Writer.Header().Get("Request-Id"),
		)

		if c.Writer.Status() >= http.StatusInternalServerError {
			entry.Error(errors.New(c.Errors.String()), "logging")
		} else {
			entry.Info(c.Request.Method + " " + c.Request.RequestURI)
		}
	}
}

// GetClientIP gets the correct IP for the end client instead of the proxy
func GetClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}

// GetDurationInMilliseconds takes a start time and returns a duration in milliseconds
func GetDurationInMilliseconds(start time.Time) time.Duration {
	end := time.Now()
	duration := end.Sub(start)
	//milliseconds := float64(duration) / float64(time.Millisecond)
	//rounded := float64(int(milliseconds*100+.5)) / 100
	return duration
}

// CorrelationID return a gin handler and add a field in the http header
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.Writer.Header().Get(CORRELATION_ID)
		if correlationID == "" {
			c.Writer.Header().Set(CORRELATION_ID, uuid.New().String())
		}
	}
}
