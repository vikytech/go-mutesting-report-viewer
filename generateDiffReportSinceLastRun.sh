!/bin/sh

curl --request GET \
    --url https://gocd.idfcbank.com/go/api/pipelines/$1/history \
    --header 'Accept: application/vnd.go.cd+json' \
    --header 'Authorization: Basic dmlraHlhdGgubl90aG86VlRvY2tlbkAxMUAy' >gocd_history

lastRunGitSha=$(cat gocd_history | jq -r ".pipelines[0].build_cause.material_revisions[0].modifications[0].revision")
lastBuildCounter=$(cat gocd_history | jq -r ".pipelines[0].counter")

curl --request GET \
    --url https://gocd.idfcbank.com/go/files/$1/$lastBuildCounter/mutation_test/1/serviceLayer/report.json \
    --header 'Accept: application/vnd.go.cd+json' \
    --header 'Authorization: Basic dmlraHlhdGgubl90aG86VlRvY2tlbkAxMUAy' >last-run-report.json

rm gocd_history
serviceFilesString=$(git diff --name-only $lastRunGitSha HEAD~1 | grep -i ".services.go")

if [ -z "$serviceFilesString" ]; then
    echo "No changes found"
    exit
fi

read -r -a serviceFiles <<<$serviceFilesString
mkdir -p reports
currentDir=$(pwd)
cd reports

for serviceFilePath in "${serviceFiles[@]}"; do
    echo "Starting to operate on $serviceFilePath"
    baseFileName=$(basename $serviceFilePath)
    mkdir -p $baseFileName
    cd $baseFileName
    go-mutesting $currentDir/$serviceFilePath &
    cd $currentDir/reports
done
wait

echo "All tests execution completed!"

cd $currentDir/reports
for serviceFilePath in "${serviceFiles[@]}"; do
    echo "Extracting report $serviceFilePath"
    baseFileName=$(basename $serviceFilePath)
    mv $baseFileName/report.json "$baseFileName.json"
    rm -r $baseFileName
done

echo "Reports extracted!"