
# Qiniu Jira Cli

`qiniujira` is a command line tool to manage Jira issues in Qiniu.

## Install

You can use the go get command install qiniujira:

```sh
go get github.com/eirture/qiniu-jira-cli/tree/main/cmd/qiniujira
```

Or clone the repo and execute the `make install` command.

## Usage

Print the help:

```sh
$ qiniujira --help                                     
qiniujira is a tool for managing jira issues

Usage:
  qiniujira [command]

Available Commands:
  completion                Generate the autocompletion script for the specified shell
  config                    Set the configs of jira and github
  help                      Help about any command
  list-deploying-issues     List all issues
  update-published-services Update published services of all associated issues

Flags:
  -h, --help   help for qiniujira

Use "qiniujira [command] --help" for more information about a command.
```

You should run the config command to add some required configs in the first time.

```sh
$ qiniujira config
Jira Address: https://jira.qiniu.io
Jira Username: <your jira username>
Jira Password: <your jira password>
You can create a new personal token at:
        https://github.com/settings/tokens/new?scopes=repo&description=qiniu-jira-cli
Github OAuth Token: <your github token>
---
Config:
jira:
    base_url: https://jira.qiniu.io
    username: <your jira username>
    password: xx************xx
github:
    oauth_token: gh************************************Cg

Is this correct? (Y/n) [Y]:
```

Now you can list deploying issues via:

```sh
$ qiniujira list-deploying-issues qboxrspub qboxs3apiv2
    ISSUE    |         MERGED AT          |  UNPUBLISHED SERVICES   
-------------+----------------------------+-------------------------
  KODO-16347 | 2023-12-14 16:17:12 +08:00 | qboxrspub, qboxs3apiv2  
  KODO-19275 | 2023-12-18 17:41:21 +08:00 | qboxs3apiv2             
  KODO-19221 | 2023-12-20 14:12:43 +08:00 | qboxs3apiv2, qboxrspub  
  KODO-19512 | 2023-12-21 16:58:51 +08:00 | qboxrspub               
  KODO-18706 | 2023-12-25 11:49:58 +08:00 | qboxs3apiv2, qboxrspub  
  KODO-19394 | 2023-12-26 10:00:35 +08:00 | qboxrspub, qboxs3apiv2  
  KODO-19300 | 2023-12-27 20:48:04 +08:00 | qboxrspub, qboxs3apiv2 
```

And you can update published services via (Only the publishing status issues will be updated):
```sh 
$ qiniujira update-published-services KODO-19447 qboxlcc
KODO-19077: []
KODO-19080: []
KODO-19083: []
KODO-19106: []
KODO-19379: []
KODO-19234: []
KODO-18966: []
KODO-17983: [qboxlcc]
KODO-18117: []
KODO-18373: []
KODO-18671: []
KODO-18824: []
KODO-18890: []
```
