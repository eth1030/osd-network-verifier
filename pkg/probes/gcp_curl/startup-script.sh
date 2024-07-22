#!/bin/sh

cat <<EOF > /usr/bin/terminate.sh
#! /bin/sh
if gcloud --quiet compute instances delete $(curl -X GET http://metadata.google.internal/computeMetadata/v1/instance/name -H 'Metadata-Flavor: Google') --zone=$(curl -X GET http://metadata.google.internal/computeMetadata/v1/instance/zone -H 'Metadata-Flavor: Google'); then : ; else
    exit 255
fi
EOF

cat <<EOF > /usr/bin/curl.sh
#! /bin/sh
if echo ${USERDATA_BEGIN} > /dev/ttyS0 ; then : ; else
    exit 255
fi
if curl --retry 3 --retry-connrefused -t B -Z -s -I -m ${TIMEOUT} -w "%{stderr}${LINE_PREFIX}%{json}\n" ${CURLOPT} ${URLS} --proto =http,https,telnet ${TLSDISABLED_URLS_RENDERED} 2>/dev/ttyS0 ; then : ; else
    exit 255
fi
if echo ${USERDATA_END} > /dev/ttyS0 ; then : ; else
    exit 255
fi
EOF

cat <<EOF > /etc/systemd/system/silence.service
[Unit]
Description=Service that silences logging to serial console

[Service]
Type=oneshot
ExecStart=systemctl mask --now serial-getty@ttyS0.service
ExecStart=systemctl disable --now syslog.socket rsyslog.service
ExecStart=sysctl -w kernel.printk="0 4 0 7"
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF


cat <<EOF > /etc/systemd/system/curl.service
[Unit]
Description=Service to run curl

[Service]
Type=oneshot
ExecStart=/usr/bin/curl.sh
Restart=on-failure
RemainAfterExit=true

[Install]
WantedBy=multi-user.target
EOF

cat <<EOF > /etc/systemd/system/terminate.service
[Unit]
Description=Service to terminate instance
[Service]
Type=oneshot
ExecStart=/usr/bin/terminate.sh
Restart=on-failure
EOF

cat <<EOF > /etc/systemd/system/terminate.timer
[Unit]
Description=Timer to terminate instance
[Timer]
OnBootSec=${DELAY}min
Unit=terminate.service
[Install]
WantedBy=multi-user.target
EOF

chmod 777 /usr/bin/curl.sh /usr/bin/terminate.sh
systemctl daemon-reload
systemctl start silence
systemctl start curl
systemctl start terminate.timer