#!/usr/bin/env bash
THRESHOLD_PERCENTAGE=95
PATH_TO_TEST=./...
ADDITIONAL_TAGS=''
PROJECT_PATH=''
CUSTOM_PATH=''
PARALLEL=''

WARNING_THRESHOLD=85
ERROR_THRESHOLD=75

FRED="\033[31m"     # foreground red
FGRN="\033[32m"     # foreground green
FYEL="\033[33m"     # foreground yellow
ENDCOLOR="\033[0m"  # end colour

for i in "$@"
do
case $i in
    -t=*|--threshold=*)
    THRESHOLD_PERCENTAGE="${i#*=}"
    shift
    ;;
    -p=*|--pathtotest=*)
    PATH_TO_TEST="./${i#*=}/..."
    shift
    ;;
    -a=*|--additionalTags=*)
    ADDITIONAL_TAGS="${i#*=}"
    shift
    ;;
    -pp=*|--projectpath=*)
    CUSTOM_PATH="${i#*=}"
    shift
    ;;
    -pl=*|--parallel=*)
    PARALLEL="-p ${i#*=}"
    shift
    ;;
    *)
    ;;
esac
done

if [ $THRESHOLD_PERCENTAGE -lt $ERROR_THRESHOLD ];then
    echo -e "$FRED Threshold is less than $ERROR_THRESHOLD !! $ENDCOLOR"
elif [ $THRESHOLD_PERCENTAGE -lt $WARNING_THRESHOLD ];then
    echo -e "$FYEL Threshold is less than $WARNING_THRESHOLD !! $ENDCOLOR"
fi

OLDPWD=$(pwd)

if [ ! -z $CUSTOM_PATH ];then
    echo -e "Testing Directory $CUSTOM_PATH"
    cd $OLDPWD/$CUSTOM_PATH
fi

EXCLUDE_PACKAGES=$(cat .coverIgnore | tr '\n' '|')

echo "Start test"

COVERAGE_WITH_PERCENTAGE=$(go test $PARALLEL -covermode=set -coverprofile $OLDPWD/coverage.out `go list $PATH_TO_TEST | grep -Ev $EXCLUDE_PACKAGES` $ADDITIONAL_TAGS | cut -d $'\t' -f '2,4' | cut -d '%' -f '1' | cut -d ':' -f '1,2')

echo "End test"

echo "Coverage for all packages :"
echo "$COVERAGE_WITH_PERCENTAGE"

cd $OLDPWD
go tool cover -html=coverage.out -o coverage.html

LOWEST_TEST_COVERAGE_PERCENTAGE=$( echo "$COVERAGE_WITH_PERCENTAGE" | grep -Eo '[0-9]+([.][0-9]+)?' | sort -n | head -1)

if [ ${LOWEST_TEST_COVERAGE_PERCENTAGE%.*} -lt $THRESHOLD_PERCENTAGE ];then
   echo "One or more packages have coverage below the threshold $THRESHOLD_PERCENTAGE !"
   exit 1
else
   echo "All packages have coverage above the threshold $THRESHOLD_PERCENTAGE !"
   exit 0
fi
