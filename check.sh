#!/bin/bash

/bin/netstat -antp | grep 3389 > /tmp/check

if [ -s /tmp/check ]
then
        rm /tmp/check
else
        cd /home/saas
        ./listener
fi