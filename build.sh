#!/bin/sh
set -e

WORKDIR=work
rm -rf $WORKDIR
mkdir -p $WORKDIR
cd $WORKDIR

echo "Downloading Loyalsoldier geosite..."
wget -q https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat -O geosite.dat

echo "Installing sing-box..."
wget -q https://github.com/SagerNet/sing-box/releases/latest/download/sing-box-linux-amd64.tar.gz
tar -xzf sing-box-linux-amd64.tar.gz
BIN=$(find . -type f -name sing-box | head -n 1)
chmod +x "$BIN"

echo "Preparing custom rules..."

TM=$(jq -R -s 'split("\n") | map(select(length>0))' ../custom/tm-rules.txt)
PROXY=$(jq -R -s 'split("\n") | map(select(length>0))' ../custom/proxy-rules.txt)

cat > custom.json <<EOF
{
  "version": 1,
  "rules": [
    {
      "domain": $TM,
      "outboundTag": "tm-rules"
    },
    {
      "domain": $PROXY,
      "outboundTag": "proxy-rules"
    }
  ]
}
EOF

echo "Merging geosite..."

$BIN geo convert \
  --input geosite.dat \
  --append custom.json \
  --output geosite-custom.dat

echo "Done"
