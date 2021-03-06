/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT

/*
Note: suggestions for level 2 subcommands doesn't work in cobra:
https://github.com/spf13/cobra/issues/981
*/

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	name            string = "uctl"
	defaultConfig   string = "/etc/" + name + ".config"
	defaultCfgEnv   string = "UCTL_CONFIG"
	defaultProject  string = ""
	defaultCertFile string = ""
)

func fatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var Verbose bool
var DryRun bool
var Cfg string = ""
var Project string = ""
var CertFile string = ""

func main() {

	var rootCmd = &cobra.Command{
		Use:   name,
		Short: name + " controls Unectio FaaS project connecting to the API endpoint",
		Long: name + ` is a CLI application to control Unectio FaaS project connecting to the API endpoint.

For online documentation, refer to https://docs.unectio.com`,
	}

	/* Completion commands and subcomands are defined here */

	var subCompletionBash = &cobra.Command{
		Use:   "completion [completion file name]",
		Short: "Generates autocompletion file for bash",
		Long: "Generates autocompletion file for bash\n" +
			"Rename it to " + name + "\n Place to /usr/local/etc/bash_completion.d \nRestart your bash",
		Args:   cobra.MinimumNArgs(1),
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletionFile(args[0])
		},
	}

	/* Function command and subcomands are defined here */

	var subFunction = &cobra.Command{
		Use:     "function",
		Aliases: []string{"fn"},
		Short:   "Manage functions and objects inside the functions (triggers, code)",
	}

	var function_enviroment string
	var subFunctionAdd = &cobra.Command{
		Use:     "add [function name]",
		Aliases: []string{"create"},
		Short:   "Add new function to the project",
		Args:    cobra.MinimumNArgs(1),
		Example: `# Create a function named "my-function"
uctl function add my-function
# Create a function named "my-function" with two environment variables
uctl function add my-function -e ENVIRONMENT=test,RUNLIMIT=35`,
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			functionAdd(&fn, &function_enviroment)
		},
	}
	subFunctionAdd.Flags().StringVarP(&function_enviroment, "environment", "e", "",
		"KEY1=VALUE1,...  set function environment variables (comma-separated for multiple)")

	var subFunctionList = &cobra.Command{
		Use:     "list [no arguments]",
		Short:   "List project functions",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			functionList()
		},
	}

	var subFunctionShow = &cobra.Command{
		Use:     "show [function name]",
		Short:   "Show function properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			info := ""
			duration := ""
			functionInfo(&fn, &info, &duration)
		},
	}

	var subFunctionDelete = &cobra.Command{
		Use:     "delete [function name]",
		Short:   "Delete the function",
		Aliases: []string{"del", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			functionDelete(&fn)
		},
	}

	var subFunctionSet = &cobra.Command{
		Use:     "set [function name]",
		Aliases: []string{"update"},
		Short:   "Modify function properties",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			functionUpdate(&fn, &function_enviroment)
		},
	}
	subFunctionSet.Flags().StringVarP(&function_enviroment, "environment", "e", "",
		"KEY1=VALUE1,...  set function environment variables (comma-separated for multiple)")

	/* Function logs command and subcomands are defined here */

	var subFunctionLogs = &cobra.Command{
		Use:   "logs [command]",
		Short: "Show function logs",
		Args:  cobra.NoArgs,
	}

	var function_log_duration string
	var subFunctionLogsShow = &cobra.Command{
		Use:   "show [function name]",
		Short: "Show function logs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			info := "logs"
			functionInfo(&fn, &info, &function_log_duration)
		},
	}
	subFunctionShow.Flags().StringVarP(&function_log_duration, "duration",
		"", "", "logs duration f.e. 100h will show logs for 100 hours back")

	/* Function stats command and subcomands are defined here */

	var subFunctionStats = &cobra.Command{
		Use:     "stats [command]",
		Short:   "Show function statistics",
		Aliases: []string{"statistics"},
		Args:    cobra.NoArgs,
	}

	var subFunctionStatsShow = &cobra.Command{
		Use:   "show [function name]",
		Short: "Show function statistics",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			info := "stats"
			duration := ""
			functionInfo(&fn, &info, &duration)
		},
	}

	/* Function code command and subcomands are defined here */

	var subFunctionCode = &cobra.Command{
		Use:   "code [command]",
		Short: "Manage function code",
		Args:  cobra.NoArgs,
	}

	var code_lang, code_src string
	var code_weight_add, code_weight_set int
	var subFunctionCodeAdd = &cobra.Command{
		Use:     "add [function name] [code name]",
		Aliases: []string{"create"},
		Short:   "Add new code of the function",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			cn := args[1]
			if code_weight_add <= 0 {
				fatal("Error: code weight must be greater than 0")
			}
			codeAdd(&fn, &cn, &code_lang, &code_src, &code_weight_add)
		},
	}
	subFunctionCodeAdd.Flags().StringVarP(&code_lang,
		"language", "l", "", "code language")
	subFunctionCodeAdd.MarkFlagRequired("language")
	subFunctionCodeAdd.Flags().StringVarP(&code_src, "source", "s", "",
		"sources (file name or url or repo:<repo name>:path)")
	subFunctionCodeAdd.MarkFlagRequired("source")
	subFunctionCodeAdd.Flags().IntVarP(&code_weight_add, "weight", "w", 1,
		"weight of the code within the function")

	var subFunctionCodeList = &cobra.Command{
		Use:     "list [function name]",
		Aliases: []string{"ls"},
		Short:   "List code of the function",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			codeList(&fn)
		},
	}

	var subFunctionCodeDelete = &cobra.Command{
		Use:     "delete [function name] [code name]",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete code of the function",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			cn := args[1]
			codeDel(&fn, &cn)
		},
	}

	var subFunctionCodeSet = &cobra.Command{
		Use:     "set [function name] [code name]",
		Aliases: []string{"update"},
		Short:   "Modify code properties",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			cn := args[1]
			if code_weight_set < 0 {
				fatal("Error: code weight must be 0 or greater")
			}

			codeSet(&fn, &cn, &code_src, &code_weight_set)
		},
	}
	subFunctionCodeSet.Flags().StringVarP(&code_src, "source", "s", "",
		"sources (file name or url or repo:<repo name>:path)")
	subFunctionCodeSet.MarkFlagRequired("source")
	subFunctionCodeSet.Flags().IntVarP(&code_weight_set, "weight", "w", 0,
		"weight of the code within the function, 0 to keep existing value")

	var subFunctionCodeShow = &cobra.Command{
		Use:     "show [function name] [code name]",
		Short:   "Show code properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			cn := args[1]
			just_code := false
			codeInfo(&fn, &cn, &just_code)
		},
	}

	var request string
	var subFunctionCodeRun = &cobra.Command{
		Use:     "run [function name] [code name]",
		Short:   "Run code of the function",
		Aliases: []string{"exec"},
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			cn := args[1]
			functionRun(&fn, &cn, &request)
		},
	}
	subFunctionCodeRun.Flags().StringVarP(&request, "request", "r", "",
		"request (JSON string)")

	/* Function code body command and subcomands are defined here */

	var subFunctionCodeBody = &cobra.Command{
		Use:     "body [command]",
		Short:   "Manage function code body",
		Aliases: []string{"source"},
		Args:    cobra.NoArgs,
	}

	var subFunctionCodeBodyShow = &cobra.Command{
		Use:   "show [function name] [code name]",
		Short: "Show function code body",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			cn := args[1]
			just_code := true
			codeInfo(&fn, &cn, &just_code)
		},
	}

	/* Function trigger command and subcomands are defined here */

	var subFunctionTrigger = &cobra.Command{
		Use:     "trigger [command]",
		Short:   "Manage function trigger",
		Aliases: []string{"tg"},
		Args:    cobra.NoArgs,
	}

	var trigger_src, trigger_url, trigger_auth, trigger_crontab, trigger_cronargs string
	var subFunctionTriggerAdd = &cobra.Command{
		Use:     "add [function name] [trigger name]",
		Aliases: []string{"create"},
		Short:   "Add new trigger of the function",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			tn := args[1]
			triggerAdd(&fn, &tn, &trigger_src, &trigger_url, &trigger_auth,
				&trigger_crontab, &trigger_cronargs)
		},
	}
	subFunctionTriggerAdd.Flags().StringVarP(&trigger_src,
		"tsource", "s", "", "trigger source (could be cron or url)")
	subFunctionTriggerAdd.Flags().StringVarP(&trigger_url, "url", "u", "",
		"trigger URL")
	subFunctionTriggerAdd.Flags().StringVarP(&trigger_auth, "auth", "a", "",
		"URL trigger auth name/id")
	subFunctionTriggerAdd.Flags().StringVarP(&trigger_crontab, "crontab", "", "",
		"cron trigger tab")
	subFunctionTriggerAdd.Flags().StringVarP(&trigger_cronargs, "cronargs", "", "",
		"Cron trigger args in foo=bar:... format")

	var subFunctionTriggerList = &cobra.Command{
		Use:     "list [function name]",
		Aliases: []string{"ls"},
		Short:   "List trigger(s) of the function",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			triggerList(&fn)
		},
	}

	var subFunctionTriggerDelete = &cobra.Command{
		Use:     "delete [function name] [trigger name]",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete trigger of the function",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			tn := args[1]
			triggerDel(&fn, &tn)
		},
	}

	var subFunctionTriggerShow = &cobra.Command{
		Use:     "show [function name] [trigger name]",
		Aliases: []string{"info"},
		Short:   "Show trigger of the function",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			tn := args[1]
			triggerInfo(&fn, &tn)
		},
	}

	/* Router command and subcomands are defined here */

	var subRouter = &cobra.Command{
		Use:     "router",
		Short:   "Manage routers",
		Aliases: []string{"rt"},
		Args:    cobra.NoArgs,
	}

	var table, table_from, rt_url string
	var subRouterAdd = &cobra.Command{
		Use:     "add [router name]",
		Aliases: []string{"create"},
		Short:   "Add new router to the project",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			routerAdd(&rn, &table, &table_from, &rt_url)
		},
	}
	subRouterAdd.Flags().StringVarP(&table, "table", "", "",
		"table (m,.../path=fn:...)")
	subRouterAdd.Flags().StringVarP(&table_from, "tablefrom", "", "",
		"file to read table from (in info -M format)")
	subRouterAdd.Flags().StringVarP(&rt_url, "url", "", "",
		"custom URL to work on")

	var subRouterSet = &cobra.Command{
		Use:     "set [router name]",
		Aliases: []string{"update"},
		Short:   "Modify router of the project",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			routerUpdate(&rn, &table, &table_from)
		},
	}
	subRouterSet.Flags().StringVarP(&table, "table", "", "",
		"table (m,.../path=fn:...)")
	subRouterSet.Flags().StringVarP(&table_from, "tablefrom", "", "",
		"file to read table from (in info -M format)")

	var subRouterList = &cobra.Command{
		Use:     "list [no arguments]",
		Short:   "List routers",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			routerList()
		},
	}

	var mux_only bool
	var subRouterShow = &cobra.Command{
		Use:     "show [router name]",
		Short:   "Show router properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			routerInfo(&rn, &mux_only)
		},
	}
	subRouterShow.Flags().BoolVarP(&mux_only, "muxonly",
		"", false, "show mux only")

	var subRouterDelete = &cobra.Command{
		Use:     "delete [router name]",
		Short:   "Delete the router",
		Aliases: []string{"del", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[0]
			routerDelete(&fn)
		},
	}

	/* Repo command and subcomands are defined here */

	var subRepo = &cobra.Command{
		Use:     "repo",
		Short:   "Manage code repositories",
		Aliases: []string{"repository"},
		Args:    cobra.NoArgs,
	}

	var rp_url string
	var subRepoAdd = &cobra.Command{
		Use:     "add [repository name]",
		Aliases: []string{"create"},
		Short:   "Add new repository to the project",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			repoAdd(&rn, &rp_url)
		},
	}
	subRepoAdd.Flags().StringVarP(&rp_url, "url", "", "",
		"repo URL (git)")
	subRepoAdd.MarkFlagRequired("url")

	var subRepoList = &cobra.Command{
		Use:     "list [no arguments]",
		Short:   "List repositories",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			repoList()
		},
	}

	var subRepoShow = &cobra.Command{
		Use:     "show [repository name]",
		Short:   "Show repository properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			repoInfo(&rn)
		},
	}

	var subRepoDelete = &cobra.Command{
		Use:     "delete [repository name]",
		Short:   "Delete the repository",
		Aliases: []string{"del", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			repoDel(&rn)
		},
	}

	var subRepoPull = &cobra.Command{
		Use:     "pull [repository name]",
		Short:   "Pull the repository",
		Aliases: []string{"download"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			repoPull(&rn)
		},
	}

	/* Repo file command and subcomands are defined here */

	var subRepoFile = &cobra.Command{
		Use:   "file [command]",
		Short: "Manage repository files",
		Args:  cobra.NoArgs,
	}

	var subRepoFileList = &cobra.Command{
		Use:     "list [repository name]",
		Aliases: []string{"ls"},
		Short:   "List files in the repository",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			fileList(&rn)
		},
	}

	var subRepoFileShow = &cobra.Command{
		Use:     "show [repository name]",
		Aliases: []string{"info"},
		Short:   "Show source of file in the repository",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rn := args[0]
			fn := args[1]
			fileInfo(&rn, &fn)
		},
	}

	/* Secret command and subcomands are defined here */

	var subSecret = &cobra.Command{
		Use:     "secret",
		Short:   "Manage secrets",
		Aliases: []string{"sc", "sec"},
		Args:    cobra.NoArgs,
	}

	var secret_key string
	var subSecretAdd = &cobra.Command{
		Use:     "add [secret name]",
		Aliases: []string{"create"},
		Short:   "Add new secret to the project",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sn := args[0]
			secretAdd(&sn, &secret_key)
		},
	}
	subSecretAdd.Flags().StringVarP(&secret_key, "key", "k", "",
		"KEY1=VALUE1,...  keys (comma-separated for multiple)")

	var subSecretList = &cobra.Command{
		Use:     "list [no arguments]",
		Short:   "List secrets",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			secretList()
		},
	}

	var subSecretShow = &cobra.Command{
		Use:     "show [secret name]",
		Short:   "Show secret properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sn := args[0]
			secretInfo(&sn)
		},
	}

	var subSecretDelete = &cobra.Command{
		Use:     "delete [secret name]",
		Short:   "Delete the secret",
		Aliases: []string{"del", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sn := args[0]
			secretDel(&sn)
		},
	}

	/* AM (authentification methods) command and subcomands are defined here */

	var subAM = &cobra.Command{
		Use:     "am",
		Short:   "Manage authentification methods",
		Aliases: []string{"auth"},
		Args:    cobra.NoArgs,
	}

	var am_key string
	var subAMAdd = &cobra.Command{
		Use:     "add [auth method name]",
		Aliases: []string{"create"},
		Short:   "Add new auth method to the project",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amn := args[0]
			authAdd(&amn, &am_key)
		},
	}
	subAMAdd.Flags().StringVarP(&secret_key, "key", "k", "",
		"jwt key (base64-encoded or auto)")

	var subAMList = &cobra.Command{
		Use:     "list [no arguments]",
		Short:   "List auth methods",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			authList()
		},
	}

	var subAMShow = &cobra.Command{
		Use:     "show [auth method name]",
		Short:   "Show auth method properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amn := args[0]
			authInfo(&amn)
		},
	}

	var subAMDelete = &cobra.Command{
		Use:     "delete [auth method name]",
		Short:   "Delete auth method",
		Aliases: []string{"del", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amn := args[0]
			authDel(&amn)
		},
	}

	/* Package command and subcomands are defined here */

	var subPackage = &cobra.Command{
		Use:     "package",
		Short:   "Manage language specific packages for functions",
		Aliases: []string{"pkg"},
		Args:    cobra.NoArgs,
	}

	var pkgl, pkgv string
	var subPackageAdd = &cobra.Command{
		Use:     "add [package name] [package language] [package version]",
		Aliases: []string{"create"},
		Short:   "Add new package to the project",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pkgn := args[0]
			packageAdd(&pkgn, &pkgl, &pkgv)
		},
	}
	subPackageAdd.Flags().StringVarP(&pkgl, "language", "l", "",
		"package language")
	subPackageAdd.MarkFlagRequired("language")
	subPackageAdd.Flags().StringVarP(&pkgv, "version", "v", "",
		"package version")

	var subPackageList = &cobra.Command{
		Use:     "list",
		Short:   "List packages with language",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			packageList(&pkgl)
		},
	}
	subPackageList.Flags().StringVarP(&pkgl, "language", "l", "",
		"package language")
	subPackageList.MarkFlagRequired("language")

	var subPackageShow = &cobra.Command{
		Use:     "show [package name]",
		Short:   "Show package properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pkgn := args[0]
			packageInfo(&pkgn, &pkgl)
		},
	}
	subPackageShow.Flags().StringVarP(&pkgl, "language", "l", "",
		"package language")
	subPackageShow.MarkFlagRequired("language")

	var subPackageDelete = &cobra.Command{
		Use:     "delete [package name]",
		Short:   "Delete package",
		Aliases: []string{"del", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pkgn := args[0]
			packageDel(&pkgn, &pkgl)
		},
	}
	subPackageDelete.Flags().StringVarP(&pkgl, "language", "l", "",
		"package language")
	subPackageDelete.MarkFlagRequired("language")

	/* Project command and subcomands are defined here */

	var subProject = &cobra.Command{
		Use:     "project",
		Short:   "Manage projects",
		Aliases: []string{"pr"},
		Args:    cobra.NoArgs,
	}

	var subProjectList = &cobra.Command{
		Use:     "list",
		Short:   "List projects",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			projectList()
		},
	}

	var subProjectShow = &cobra.Command{
		Use:     "show [project id]",
		Short:   "Show project properties",
		Aliases: []string{"info"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			prn := args[0]
			projectInfo(&prn)
		},
	}

	/* Global flags definition */

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "", false,
		"verbose output (print raw REST API data)")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry", "", false,
		"perform all client-side validation but do not perform any server requests")
	rootCmd.PersistentFlags().StringVarP(&Cfg, "config", "", "",
		"path to the configuration file")
	rootCmd.PersistentFlags().StringVarP(&Project, "project", "", defaultProject,
		"project id")

	/* To not forget howto generate certificate
	   	openssl req \
	       -x509 \
	       -nodes \
	       -newkey rsa:2048 \
	       -keyout server.key \
	       -out server.crt \
	       -days 3650 \
	       -subj "/C=GB/ST=London/L=London/O=Global Security/OU=IT Department/CN=*"
	*/
	rootCmd.PersistentFlags().StringVarP(&CertFile, "certificate", "", defaultCertFile,
		"path to server certificate file for https connection.\n"+
			"By default systems certificates pool will be used.\n"+
			"On windows you need to use this option always")

	/* CLI commands initialisation */

	rootCmd.AddCommand(subCompletionBash)
	rootCmd.AddCommand(subFunction)
	rootCmd.AddCommand(subRouter)
	rootCmd.AddCommand(subRepo)
	rootCmd.AddCommand(subSecret)
	rootCmd.AddCommand(subAM)
	rootCmd.AddCommand(subPackage)
	rootCmd.AddCommand(subProject)

	subFunction.AddCommand(subFunctionAdd)
	subFunction.AddCommand(subFunctionList)
	subFunction.AddCommand(subFunctionShow)
	subFunction.AddCommand(subFunctionLogs)
	subFunction.AddCommand(subFunctionStats)
	subFunction.AddCommand(subFunctionDelete)
	subFunction.AddCommand(subFunctionSet)
	subFunction.AddCommand(subFunctionCode)
	subFunction.AddCommand(subFunctionTrigger)

	subFunctionLogs.AddCommand(subFunctionLogsShow)

	subFunctionStats.AddCommand(subFunctionStatsShow)

	subFunctionCode.AddCommand(subFunctionCodeAdd)
	subFunctionCode.AddCommand(subFunctionCodeList)
	subFunctionCode.AddCommand(subFunctionCodeDelete)
	subFunctionCode.AddCommand(subFunctionCodeSet)
	subFunctionCode.AddCommand(subFunctionCodeShow)
	subFunctionCode.AddCommand(subFunctionCodeRun)

	subFunctionCode.AddCommand(subFunctionCodeBody)

	subFunctionCodeBody.AddCommand(subFunctionCodeBodyShow)

	subFunctionTrigger.AddCommand(subFunctionTriggerAdd)
	subFunctionTrigger.AddCommand(subFunctionTriggerList)
	subFunctionTrigger.AddCommand(subFunctionTriggerDelete)
	subFunctionTrigger.AddCommand(subFunctionTriggerShow)

	subRouter.AddCommand(subRouterAdd)
	subRouter.AddCommand(subRouterSet)
	subRouter.AddCommand(subRouterList)
	subRouter.AddCommand(subRouterDelete)
	subRouter.AddCommand(subRouterShow)

	subRepo.AddCommand(subRepoAdd)
	subRepo.AddCommand(subRepoList)
	subRepo.AddCommand(subRepoDelete)
	subRepo.AddCommand(subRepoShow)
	subRepo.AddCommand(subRepoPull)

	subRepo.AddCommand(subRepoFile)
	subRepoFile.AddCommand(subRepoFileList, subRepoFileShow)

	subSecret.AddCommand(subSecretAdd)
	subSecret.AddCommand(subSecretList)
	subSecret.AddCommand(subSecretDelete)
	subSecret.AddCommand(subSecretShow)

	subAM.AddCommand(subAMAdd)
	subAM.AddCommand(subAMList)
	subAM.AddCommand(subAMDelete)
	subAM.AddCommand(subAMShow)

	subPackage.AddCommand(subPackageAdd)
	subPackage.AddCommand(subPackageList)
	subPackage.AddCommand(subPackageDelete)
	subPackage.AddCommand(subPackageShow)

	subProject.AddCommand(subProjectList, subProjectShow)

	rootCmd.Execute()
}
