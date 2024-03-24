# Copy service file, incase if there are any changes
cd ~/gocron
sudo cp bin/* /gocron/
sudo cp gocron-web.service /etc/systemd/system/gocron-web.service
sudo cp gocron-node.service /etc/systemd/system/gocron-node.service
# reload configurations incase if service file has changed
sudo systemctl daemon-reload
# restart the service
sudo systemctl restart gocron-web
sudo systemctl restart gocron-node
# start of VM restart
sudo systemctl enable gocron-web
sudo systemctl enable gocron-node