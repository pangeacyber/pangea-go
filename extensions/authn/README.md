# AuthN Extension to provide user import functionality

This is AuthN extension utility which provides the functionality to import one time users to Pangea from given CSV or
Json File

## Description

This is command line utility which allow customer to import users from given CSV or Json file

## Installation
* Install go lang 1.20 onwards
* Run `go install`
* you will see `authn` command in `$GOROOT/bin` folder

## Command line
```bash
 ~/go/bin/authn import -h
One time user import to the pangea.

Usage:
  authn import [filePath] [flags]

Flags:
  -f, --dry-run                mimic run import workflow (it does not make api call to create users). Default is false
  -m, --fieldsMapping string   Fields mapping file to map source provider to pangea
  -h, --help                   help for import
  -i, --importFile string      import user csv or json file

Global Flags:
  -c, --config string   config file (default is $HOME/.pangea.yaml)
  -d, --domain string   pangea domain (default is PANGEA_DOMAIN env)
  -t, --token string    pangea token (default is PANGEA_TOKEN env)
```
