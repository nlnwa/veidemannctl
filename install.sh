#!/usr/bin/env bash

VERSION=`curl -s -I https://github.com/nlnwa/veidemannctl/releases/latest | grep Location | sed -e 's/.*tag\/\([0-9\.]\+\).*$/\1/'`
echo Installing veidemannctl v${VERSION}

curl -Lo veidemannctl https://github.com/nlnwa/veidemannctl/releases/download/${VERSION}/veidemannctl_${VERSION}_linux_amd64
sudo mv veidemannctl /usr/local/bin/veidemannctl
sudo chmod +x /usr/local/bin/veidemannctl

# Install command completion
sudo sh -c "/usr/local/bin/veidemannctl completion > /etc/bash_completion.d/veidemannctl"
