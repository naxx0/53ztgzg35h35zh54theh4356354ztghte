#!/bin/bash

/bin/netstat -antp | grep 20211 > /tmp/check

if [ -s /tmp/check ]
then
        rm /tmp/check
else
        cd /opt/stacks/pialert
        docker compose up -d
fi