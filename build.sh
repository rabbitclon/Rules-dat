#!/bin/sh
set -e

WORKDIR=work
XRAYDIR=xray

rm -rf $WORKDIR
rm -rf $XRAYDIR

mkdir -p $WORKDIR
mkdir -p $XRAYDIR

cd $WORKDIR

echo "Download Loyalsoldier geosite..."
wget -q https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat -O geosite.dat

cd ..

echo "Download Xray core..."
wget -q https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip -O xray.zip

unzip -oq xray.zip -d $XRAYDIR

XRAY=$(find $XRAYDIR -type f -name xray | head -n 1)
chmod +x "$XRAY"

echo "Prepare custom rules..."

TM=$(grep -v '^$' custom/tm-rules.txt | sed 's/^/"/;s/$/"/' | paste -sd, -)
PROXY=$(grep -v '^$' custom/proxy-rules.txt | sed 's/^/"/;s/$/"/' | paste -sd, -)

cat > custom.json <<EOF
{
  "version": 1,
  "rules": [
    {
      "domain": [$TM],
      "outboundTag": "tm-rules"
    },
    {
      "domain": [$PROXY],
      "outboundTag": "proxy-rules"
    }
  ]
}
EOF

echo "Build geosite..."

$XRAY run geo \
  --input-file=work/geosite.dat \
  --domain-file=custom.json \
  --output=work/geosite-custom.dat

echo "Done!"
