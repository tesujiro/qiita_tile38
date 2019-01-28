#!/bin/bash

commands-1(){
        cat <<EOF | awk 'gsub(/#.*/,"")>=0'
DROP location
DROP lunch
#
GET location me POINT 35.6581 139.6975
# 
SET lunch ramen:A POINT 35.6586 139.6982
SET lunch gyudon:B POINT 35.6570 139.6967
SET lunch ramen:C POINT 35.6576 139.6948
SCAN lunch MATCH ramen* IDS
NEARBY lunch POINT 35.6581 139.6975 100
EOF
}

commands-2(){
        cat <<EOF | awk 'gsub(/#.*/,"")>=0'
DROP location
DROP example
#
SET location me POINT 35.6581 139.6975
SET example bounds:X BOUNDS 35.6578 139.6971 35.6581 139.6968
SET example bounds:Y BOUNDS 35.6572 139.6984 35.6575 139.6978
SET example bounds:Z BOUNDS 35.6590 139.6967 35.6594 139.6959
WITHIN example IDS CIRCLE 35.6581 139.6975 120
INTERSECTS example IDS CIRCLE 35.6581 139.6975 120
EOF
}
#commands-2 | tile38-cli

geoJson(){
	local TYPE=$1
	local COORDINATES=$2

	cat <<EOF | awk 'gsub(/#.*/,"")>=0' | tr -d '\t' | tr -d '\n'
{
		"type":"$TYPE",
		"coordinates":$COORDINATES
	}
EOF
}

PROMPT="127.0.0.1:9851> "
exec_tile38_cli(){
	echo -n "$PROMPT"
	echo "$@"
	echo "$@" | tile38-cli
}
export -f exec_tile38_cli

commands-3(){
		local KEY=example

        cat <<EOF | awk 'gsub(/#.*/,"")>=0'
DROP location
DROP $KEY
SET location me OBJECT $(geoJson Polygon [[[35.6590,139.6982],[35.6589,139.6978],[35.6577,139.6965],[35.6574,139.6964],[35.6572,139.6966],[35.6575,139.6973],[35.6580,139.6988],[35.6587,139.6984],[35.6590,139.6982]]])
SET $KEY polygon:P OBJECT $(geoJson Polygon [[[35.6587,139.6984],[35.6590,139.6983],[35.6589,139.6979],[35.6586,139.6980],[35.6587,139.6984]]])
SET $KEY polygon:Q OBJECT $(geoJson Polygon [[[35.6591,139.6967],[35.6595,139.6960],[35.6589,139.6958],[35.6586,139.6965],[35.6591,139.6967]]])
SET $KEY road:R OBJECT $(geoJson LineString [[35.6584,139.6954],[35.6567,139.6970]])
SET $KEY road:S OBJECT $(geoJson LineString [[35.6585,139.6994],[35.6575,139.6953]])
INTERSECTS $KEY IDS GET location me
EOF
}

commands-3 | while read line
do
	exec_tile38_cli "$line"
done

