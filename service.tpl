[Unit]
Description=A bridge between 433Mhz switches and Philips Hue
After=network.target

[Service]
Type=simple
User=__DAEMON_USER__
Group=__DAEMON_USER__
WorkingDirectory=/home/__DAEMON_USER__
ExecStart=/home/__DAEMON_USER__/hue-lightswitch -hue-host=__HUE_HOST__
Restart=always

[Install]
WantedBy=default.target
