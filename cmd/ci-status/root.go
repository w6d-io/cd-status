/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 05/05/2023
*/

package ci_status

import (
    "github.com/w6d-io/ci-status/cmd/ci-status/server"
    "os"

    "github.com/spf13/cobra"

    "github.com/w6d-io/ci-status/internal/config"
    "github.com/w6d-io/x/cmdx"
    "github.com/w6d-io/x/logx"
)

var rootCmd = &cobra.Command{
    Use: "ci-status",
}

var OsExit = os.Exit

func Execute() {
    log := logx.WithName(nil, "Main.Command")
    //info, _ := debug.ReadBuildInfo()
    //fmt.Println(info)
    rootCmd.AddCommand(cmdx.Version(&config.Version, &config.Revision, &config.Built))
    rootCmd.AddCommand(server.Cmd)
    if err := rootCmd.Execute(); err != nil {
        log.Error(err, "exec command failed")
        OsExit(1)
    }
}
