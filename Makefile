tidy:
	go mod tidy

build:
	go build -o medialog

install:
	sudo systemctl stop medialog
	sudo mkdir /usr/lib/medialog/prev
	sudo mv /usr/lib/medialog/medialog /usr/lib/medialog/prev
	sudo mv /usr/lib/medialog/public /usr/lib/medialog/prev
	sudo mv /usr/lib/medialog/templates /usr/lib/medialog/prev
	sudo cp medialog /usr/lib/medialog/
	sudo cp -r templates /usr/lib/medialog
	sudo cp -r public /usr/lib/medialog
	sudo chown -R centos:centos /usr/lib/medialog

clean:
	rm medialog