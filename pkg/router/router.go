/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 23/01/2021
*/

package router

import (
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/gin-gonic/gin"
    "github.com/w6d-io/ci-status/internal/config"
)

var Address string

func init() {
    engine.Use(LogMiddleware())
    engine.Use(gin.Recovery())
    if config.IsAuth() {
        engine.Use(Auth())
    }
}

// AddPOST adds handler and path to the engine
func AddPOST(relativePath string, handlers ...gin.HandlerFunc) {
    engine.POST(relativePath, handlers...)
}

// AddGET adds handler and path to the engine
func AddGET(relativePath string, handlers ...gin.HandlerFunc) {
    engine.GET(relativePath, handlers...)
}

// Run execute le gin router
func Run() error {

    server.Addr = Address
    server.Handler = engine

    quit := make(chan os.Signal)
    signal.Notify(quit, syscall.SIGINT)
    signal.Notify(quit, syscall.SIGKILL)
    go func() {
        <-quit
        logger.Info("receive interrupt or kill signal")
        if err := server.Close(); err != nil {
            logger.Error(err, "Server closed")
            os.Exit(1)
        }
    }()
    logger.WithValues("address", Address).Info("Listening and serving HTTP")
    if err := server.ListenAndServe(); err != nil {
        if err == http.ErrServerClosed {
            logger.Info("Server closed under request")
            return nil
        }
        logger.Error(err, "Server closed unexpect")
        return err
    }
    return nil
}

// Stop the http server
func Stop() error {
    if server != nil {
        return server.Close()
    }
    return nil
}
