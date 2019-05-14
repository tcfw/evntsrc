#!/bin/bash -e

detect_changed_services() {
	echo "----------------------------------------------"
 	echo "detecting changed folders for this commit"

	# get a list of all the changed folders only
	changed_folders=`git diff --name-only $COMMIT_RANGE | grep internal/ | awk 'BEGIN {FS="/"} {print $2}' | uniq`
	echo "changed folders "$changed_folders

	changed_services=()
	for folder in $changed_folders
	do
		if [ "$folder" == 'utils' ] || [ "$folder" == 'tracing' ] || [ "$folder" == 'event' ]; then
			echo "!! a common folder changed, building and publishing all microservices"
			changed_services=`find ./internal -maxdepth 1 -type d -not -name 'utils' -not -name 'tracing' -not -name 'event' -not -name '.git' -not -path './internal' | sed 's|./internal/||'`
			break
		else
			echo "Adding $folder to list of services to build"
			changed_services+=("$folder")
		fi
	done

	#update APIGW if any proto files change
	proto_change_count=`git diff --name-only $COMMIT_RANGE | grep "\.proto$" | wc -l`
	if [ $proto_change_count -gt 0 ]; then
		changed_services+=("apigw")
		echo "Added APIGW due to proto change"
	fi

	changed_services=( $(printf '%s\n' "${changed_services[@]}" | sort -u) )

	# Iterate on each service and run the packaging script
	for service in ${changed_services[@]}
	do
		echo ""
		echo "-------------------Running packaging for $service---------------------"
		make $service
	done
}

detect_changed_services