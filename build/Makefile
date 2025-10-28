v := $(shell $(MEDIALOG_HOME)/medialog --version | jq '.version')
ts := $(shell date +%s)
ver := $(v)-$(ts)

upgrade: stop archive install start

archive:
	sudo mkdir $(MEDIALOG_HOME)/$(ver)
	sudo mv $(MEDIALOG_HOME)/medialog $(MEDIALOG_HOME)/$(ver)
	sudo mv $(MEDIALOG_HOME)/public $(MEDIALOG_HOME)/$(ver)
	sudo mv $(MEDIALOG_HOME)/templates $(MEDIALOG_HOME)/$(ver)
	sudo tar cvzf $(MEDIALOG_HOME)/$(ver).tgz $(MEDIALOG_HOME)/$(ver)
	sudo mv $(MEDIALOG_HOME)/$(ver).tgz $(MEDIALOG_HOME)/previous-versions/
	sudo rm -r $(MEDIALOG_HOME)/$(ver)
	sudo chown -R medialog:medialog $(MEDIALOG_HOME)

install:
	chmod +x medialog
	sudo cp medialog $(MEDIALOG_HOME)
	sudo cp -r templates/ public/ $(MEDIALOG_HOME)
	sudo chown -R medialog:medialog $(MEDIALOG_HOME)

deploy:
	@echo "Not Implemented Yet"

update-templates:
	sudo cp -r templates $(MEDIALOG_HOME)
	sudo chown -R medialog:medialog $(MEDIALOG_HOME)

version:
	@echo $(v)

start:
	sudo systemctl start medialog

stop:
	sudo systemctl stop medialog