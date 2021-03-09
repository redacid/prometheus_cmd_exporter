# Install
go build
mkdir /usr/sbin/cmd-exporter 
cp prometheus_cmd_metrics /usr/sbin/cmd-exporter/cmd-exporter
cp config.json /usr/sbin/cmd-exporter/

/usr/lib/systemd/system/nginx_exporter.service



