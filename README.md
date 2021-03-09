# Install
go build
mkdir /usr/sbin/cmd-exporter 
cp prometheus_cmd_metrics /usr/sbin/cmd-exporter/cmd-exporter
cp config.json /usr/sbin/cmd-exporter/
cp cmd-exporter.service /usr/lib/systemd/system/cmd-exporter.service



