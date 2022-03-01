#!/bin/bash
cd /home/ubuntu/
sudo snap install go --classic
git clone https://github.com/Botiyava/DelftSecure
mkdir  /home/ubuntu/DelftSecure/log
mkdir  /home/ubuntu/DelftSecure/configs
touch /home/ubuntu/DelftSecure/log/errors.log
chmod 666 /home/ubuntu/DelftSecure/log/errors.log


cat <<'EOF' > /home/ubuntu/DelftSecure/configs/config1.json
${config1}
EOF