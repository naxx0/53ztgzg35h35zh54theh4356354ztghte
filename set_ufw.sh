#!/bin/bash

####### Setting ufw ports for listener

ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 8080/tcp
ufw allow 21/tcp
ufw allow 23/tcp
ufw allow 25/tcp
ufw allow 9090/tcp
ufw allow 3389/tcp
ufw allow 135/tcp
ufw allow 139/tcp
ufw allow 445/tcp
ufw allow 3306/tcp
ufw allow 389/tcp
ufw allow 7070/tcp
ufw allow 8443/tcp
ufw allow 1433/tcp

ufw allow 22/tcp

yes | ufw enable