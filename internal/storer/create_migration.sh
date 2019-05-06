#!/bin/bash

if [ "$1" = "" ]; then
	echo "No migration name given";
	exit 0;
fi

ds=$(date +%s)
touch "./migrations/${ds}_$1.sql"