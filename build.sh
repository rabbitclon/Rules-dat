#!/bin/sh
set -e

WORKDIR=work
rm -rf $WORKDIR
mkdir -p $WORKDIR
cd $WORKDIR

echo "Download Loyalsoldier geosite..."
wget -q https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat -O geosite.dat

echo "Download Xray core..."
wget -q https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip
unzip -q Xray-linux-64.zip

XRAY=$(find . -type f -name xray | head -n 1)
chmod +x "$XRAY"

echo "Preparing custom rules..."

TM=$(grep -v '^$' ../custom/tm-rules.txt | sed 's/^/"/;s/$/"/' | paste -sd, -)
PROXY=$(grep -v '^$' ../custom/proxy-rules.txt | sed 's/^/"/;s/$/"/' | paste -sd, -)

cat > geosite_custom.json <<EOF
{
  "version": 1,
  "rules": [
    {
      "domain": [$TM],
      "outboundTag": "tm-rules"
    },
    {
      "domain": [$PROXY],
      "outboundTag": "proxy"
    }
  ]
}
EOF

echo "Merging via Xray geo..."

$XRAY run geo \
  --input-file=geosite.dat \
  --domain-file=geosite_custom.json \
  --output=geosite-custom.dat

echo "Done"
