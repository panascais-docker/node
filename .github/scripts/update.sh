#!/bin/bash

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd );
url="https://unofficial-builds.nodejs.org/download/release";
versions="";
checksums="";

declare -a mapping;
mapping['4']="https://nodejs.org/dist/latest-v4.x/";
mapping['5']="https://nodejs.org/dist/latest-v5.x/";
mapping['6']="https://nodejs.org/dist/latest-v6.x/";
mapping['7']="https://nodejs.org/dist/latest-v7.x/";
mapping['8']="https://nodejs.org/dist/latest-v8.x/";
mapping['9']="https://nodejs.org/dist/latest-v9.x/";
mapping['10']="https://nodejs.org/dist/latest-v10.x/";
mapping['11']="https://nodejs.org/dist/latest-v11.x/";
mapping['12']="https://nodejs.org/dist/latest-v12.x/";
mapping['13']="https://nodejs.org/dist/latest-v13.x/";
mapping['14']="https://nodejs.org/dist/latest-v14.x/";
mapping['15']="https://nodejs.org/dist/latest-v15.x/";
mapping['16']="https://nodejs.org/dist/latest-v16.x/";
# manually maintain this and replace it with v16.x on 2021-10-26
mapping['98']="https://nodejs.org/dist/latest-v14.x/";
mapping['99']="https://nodejs.org/dist/latest/";

for key in ${!mapping[@]}; do
    val=${mapping[$key]};
    raw="$( curl -sSL $val | grep node-v | head -1; )";
    IFS='-' read -ra lines <<< "$raw"
    for line in "${lines[@]}"; do
        if [[ $line = v* ]]; then
            result=$( echo $line | cut -c 2-$((${#line})) );
            name=$( if [[ $key -eq '98' ]]; then echo 'lts'; elif [[ $key -eq '99' ]]; then echo 'latest'; else echo $key; fi; );
            hash=$( curl -sSL $url/$line/SHASUMS256.txt | grep musl.tar.xz | head -1 | cut -d ' ' -f1 | tr '\n' ' '; );
            versions+="$name=$result ";
            if [[ -n $hash ]]; then checksums+="$name=$hash "; fi;
            echo "Completed update for version=$name";
            break;
        fi;
    done;
    unset IFS;
done;

jo -p $versions > $dir/../../configuration/versions.json;
jo -p $checksums > $dir/../../configuration/checksums.json;

# Since jo doesn't pretty print JSON as we would like it to do, 
# we just use a hacky node script to reformat the json after it
# is written.
node $dir/update.js versions;
node $dir/update.js checksums;

echo "Configuration data written to disk.";