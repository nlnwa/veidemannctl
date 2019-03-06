#!/usr/bin/env bash

VERSION=`curl -s -I https://github.com/nlnwa/veidemannctl/releases/latest | grep Location | sed -e 's/.*tag\/\([0-9\.]\+\).*$/\1/'`
echo Installing veidemannctl v${VERSION}

wget -q --show-progress https://github.com/nlnwa/veidemannctl/releases/download/${VERSION}/veidemannctl_${VERSION}_linux_amd64
sudo cp veidemannctl_${VERSION}_linux_amd64 /usr/local/bin/veidemannctl
sudo chmod +x /usr/local/bin/veidemannctl
rm veidemannctl_${VERSION}_linux_amd64
sudo sh -c "veidemannctl completion > /etc/bash_completion.d/veidemannctl"
