#!/usr/bin/env bash

VERSION=$(curl -s -I https://github.com/nlnwa/veidemannctl/releases/latest | grep location | sed -e 's/.*tag\/\([0-9\.]\+\).*$/\1/')
echo "Installing veidemannctl v${VERSION}"

curl -Lo veidemannctl https://github.com/nlnwa/veidemannctl/releases/download/${VERSION}/veidemannctl_${VERSION}_linux_amd64
sudo install veidemannctl /usr/local/bin/veidemannctl
rm veidemannctl

# Install command completion for bash
sudo sh -c "/usr/local/bin/veidemannctl completion bash > /etc/bash_completion.d/veidemannctl"
