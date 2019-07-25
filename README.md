# PBX Getbox renders download app

This script checks and downloads render from printbox server to local server.

## Usage

If you have docker installed in server then use ```docker-compose up -d```
This runs app forever and checking. App is doing interval checks after time set in config.

Running only once you need server to have go installed and then use ```go run main.go```

## Config setup

```yml
host: qa-fotomeened-pbx2.getprintbox.com
schema: https
api_url: api/getbox
api_version: v2
auth_access_key: access_key_here
auth_secret_key: secret_key_here
getbox_id: getbox id here #ask printbox
nr_of_cores: 8 #number of cores/fetchers you would like to use
getbox_path: /var/www/getbox/ #location of getbox files root
temp_path: tmp/ #temporary files dir. Unzipped versions of renders will be downloaded here
finished_path: finished/ #unzipped versions of files will be placed here
fetchers_interval: 10 #seconds interval you would like getbox app run fetchers
```

## Development

Run tests ```go test ./...```