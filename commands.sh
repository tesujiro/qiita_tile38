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
commands-2 | tile38-cli
