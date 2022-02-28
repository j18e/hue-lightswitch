[Unit]
Description=A bridge between 433Mhz switches and Philips Hue
After=network.target

[Service]
Type=simple
User=pi
Group=pi
WorkingDirectory=/home/pi
ExecStart=/home/pi/hue-lightswitch -hue-host=__HUE_HOST__
Restart=always

[Install]
WantedBy=default.target
