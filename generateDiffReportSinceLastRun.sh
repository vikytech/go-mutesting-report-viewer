#!/usr/bin/env bash

curl --request GET \
    --url https://gocd.idfcbank.com/go/api/pipelines/$1/history \
    --header 'Accept: application/vnd.go.cd+json' \
    --header 'Authorization: Basic dmlraHlhdGgubl90aG86VlRvY2tlbkAxMUAy' >gocd_history

lastRunGitSha=$(cat gocd_history | jq -r ".pipelines[0].build_cause.material_revisions[0].modifications[0].revision")
lastBuildCounter=$(cat gocd_history | jq -r ".pipelines[0].counter")

curl --request GET \
    --url https://gocd.idfcbank.com/go/files/$1/$lastBuildCounter/mutation_test/1/serviceLayer/report.json \
    --header 'Accept: application/vnd.go.cd+json' \
    --header 'Authorization: Basic dmlraHlhdGgubl90aG86VlRvY2tlbkAxMUAy' >last_run_report.json

rm gocd_history

servicefiles=$(git diff --name-only $lastRunGitSha HEAD~1 | grep -i ".service.go")

if [ -z "$servicefiles" ]; then
    echo "No changes found"
    exit
fi

go-mutesting $servicefiles
# ... trigger diff plugin and upload new report
rm last_run_report.json