tidy:
	go mod tidy

build:
	go build -o medialog

archive:
	sudo systemctl stop medialog
	sudo mkdir $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/medialog $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/public $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/templates $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/prev $(MEDIALOG_HOME)/previouse-versions
	sudo chown -R centos:centos $(MEDIALOG_HOME)

install:
	sudo systemctl stop medialog
	sudo cp medialog $(MEDIALOG_HOME)/
	sudo cp -r templates /$(MEDIALOG_HOME)/
	sudo cp -r public $(MEDIALOG_HOME)/
	sudo chown -R centos:centos $(MEDIALOG_HOME)

clean:
	rm medialog

update-templates:
	sudo cp -r templates /$(MEDIALOG_HOME)
	sudo chown -R centos:centos /$(MEDIALOG_HOME)