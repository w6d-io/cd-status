/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 07/05/2023
*/

package config

import (
    "context"
    "github.com/spf13/viper"
    "github.com/w6d-io/hook"
    "github.com/w6d-io/x/errorx"
    "github.com/w6d-io/x/logx"
)

type Hook struct {
    URL   string `json:"url"`
    Scope string `json:"scope"`
}

func hookSubscription() error {
    var hooks []Hook
    log := logx.WithName(nil, "Hook.Subscription")
    if err := viper.UnmarshalKey(ViperKeyHooks, &hooks); err != nil {
        log.Error(err, "unmarshalling hook failed")
        return errorx.New(err, "unmarshalling hook failed")
    }
    log.V(2).Info("subscripting", "count", len(hooks))
    for _, h := range hooks {
        if err := hook.Subscribe(context.Background(), h.URL, h.Scope); err != nil {
            log.Error(err, "subscription failed", "url", h.URL, "scope", h.Scope)
            return errorx.New(err, "subscription failed")
        }
        log.V(2).Info("subscription", "url", h.URL, "scope", h.Scope)
    }
    return nil
}
