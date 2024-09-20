#!/usr/bin/env bash

curl --request GET \
  --url https://gocd.idfcbank.com/go/api/pipelines/$1/history \
  --header 'Accept: application/vnd.go.cd+json' \
  --header 'Authorization: Basic dmlraHlhdGgubl90aG86VlRvY2tlbkAxMUAy' \
  --header 'User-Agent: insomnia/8.4.5' \
  --cookie JSESSIONID=node058t3hy0jr4121lek1tffz3m8r4290205.node0 > gocd_history

# echo $history
lastRunGitSha=`cat gocd_history | jq -r ".pipelines[0].build_cause.material_revisions[0].modifications[0].revision"`
rm gocd_history
echo $lastRunGitSha

servicefiles=$(git diff --name-only $lastRunGitSha HEAD~1 | grep -i ".service.go")

[ -z "$servicefiles" ] && echo "No changes found" || go-mutesting $servicefiles
