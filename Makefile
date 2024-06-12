tidy:
	go mod tidy

clean:
	rm bin/medialog

build:
	go build -o bin/medialog medialog.go

archive:
	sudo systemctl stop medialog
	sudo mkdir $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/medialog $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/public $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/templates $(MEDIALOG_HOME)/prev
	sudo mv $(MEDIALOG_HOME)/prev $(MEDIALOG_HOME)/previous-versions
	sudo chown -R medialog:medialog $(MEDIALOG_HOME)

install:
	sudo systemctl stop medialog
	sudo cp bin/medialog $(MEDIALOG_HOME)/
	sudo cp -r templates $(MEDIALOG_HOME)/
	sudo cp -r public $(MEDIALOG_HOME)/
	sudo chown -R medialog:medialog $(MEDIALOG_HOME)

update-templates:
	sudo cp -r templates $(MEDIALOG_HOME)
	sudo chown -R medialog:medialog $(MEDIALOG_HOME)