tidy:
	go mod tidy

build:
	go build -o medialog

archive:
	sudo systemctl stop medialog
	sudo mkdir /var/www/medialog/prev
	sudo mv /var/www/medialog/medialog /var/www/medialog/prev
	sudo mv /var/www/medialog/public /var/www/medialog/prev
	sudo mv /var/www/medialog/templates /var/www/medialog/prev
	sudo mv /var/www/medialog/prev /var/www/medialog/previouse-versions
	sudo chown -R centos:centos /var/www/medialog

install:
	sudo systemctl stop medialog
	sudo cp medialog /var/www/medialog/
	sudo cp -r templates /var/www/medialog
	sudo cp -r public /var/www/medialog
	sudo chown -R centos:centos /var/www/medialog

clean:
	rm medialog

update-templates:
	sudo cp -r templates /var/www/medialog
	sudo chown -R centos:centos /var/www/medialog