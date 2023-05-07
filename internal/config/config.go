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
Created on 22/01/2021
*/

package config

import (
    "context"
    "github.com/fsnotify/fsnotify"
    "github.com/w6d-io/ci-status/internal/tekton"
    "github.com/w6d-io/ci-status/pkg/router"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/spf13/viper"
    "github.com/w6d-io/ci-status/internal/embedx"
    "github.com/w6d-io/jsonschema"
    "github.com/w6d-io/x/cmdx"
    "github.com/w6d-io/x/logx"
)

var (
    // Version of application
    Version string

    // Revision is the commit of this version
    Revision string

    // Built is the timestamp od this version
    Built string

    // CfgFile contain the path of the config file
    CfgFile string

    // OsExit is hack for unit-test
    OsExit = os.Exit

    // SkipValidation toggling the config validation
    SkipValidation bool
)

func setDefault() {
    viper.SetDefault(ViperKeyHTTPListen, ":8080")
    viper.SetDefault(ViperTimeout, time.Hour)
}

func Init() {
    base := filepath.Base(CfgFile)
    log := logx.WithName(nil, "Config.Init")
    ext := filepath.Ext(CfgFile)
    log.V(2).Info("viper",
        "path", CfgFile,
        "ext", filepath.Ext(CfgFile),
        "type", strings.TrimLeft(ext, "."),
        "configName", FileNameWithoutExtension(base),
        "base", base,
        "dir", filepath.Dir(CfgFile),
    )
    setDefault()
    viper.SetConfigName(FileNameWithoutExtension(base))
    viper.SetConfigType(strings.TrimLeft(ext, "."))
    viper.AddConfigPath(filepath.Dir(CfgFile))
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME/.ci_status")

    viper.SetEnvPrefix("cs")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            logx.WithName(context.TODO(), "Config").Error(err, "failed to read config")
            OsExit(1)
            return
        }
        log.Error(err, "config not found")
        return
    }
    var c Config
    cmdx.Should(viper.Unmarshal(&c), "unmarshal config failed")
    if !SkipValidation {
        log.Info("run config validation")
        cmdx.Must(jsonschema.AddSchema(jsonschema.Config, embedx.ConfigSchema), "add config schema failed")
        cmdx.Must(jsonschema.Config.Validate(&c), "config validation failed")
    }
    log.Info("config loaded", "file", viper.ConfigFileUsed())
    cmdx.Must(hookSubscription(), "hook subscription failed")
    initClient()
    viper.WatchConfig()
    viper.OnConfigChange(onChange())
    return
}

func initClient() {
    router.Address = viper.GetString(ViperKeyHTTPListen)
    tekton.Timeout = viper.GetDuration(ViperTimeout)
    if config.Timeout == 0 {
        config.Timeout = viper.GetDuration(ViperTimeout)
    }
}
func onChange() func(event fsnotify.Event) {
    return func(e fsnotify.Event) {
        log := logx.WithName(nil, "Config").WithName("onChange")
        log.Info("config changed", "file", e.Name)

        initClient()
    }
}

// FileNameWithoutExtension returns the
func FileNameWithoutExtension(fileName string) string {
    return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
