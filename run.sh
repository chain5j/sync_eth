#!/bin/bash

Version=$1
if [ -z "$Version" ]; then
    echo "Version is empty"
    exit 25
fi


docker rm -f sync_eth
docker run -itd --restart=unless-stopped -v /etc/localtime:/etc/localtime -v /etc/timezone:/etc/timezone --name sync_eth -v `pwd`/project:/data --network=host sync_eth:${Version}
docker logs -f sync_eth
