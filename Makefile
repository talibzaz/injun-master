all:
	make build-binary
	make build
	make run
	rm ./injun-linux-amd64

build-binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o injun-linux-amd64 .

build:
	docker build -t gw/injun:latest .

bunldle-assets:
	rice clean && rice embed-go

run:
	$(eval ip=$(shell docker inspect --format '{{ .NetworkSettings.Networks.bridge.IPAddress }}' gelffy))\
	docker run  -d \
		--log-driver=fluentd \
		--log-opt fluentd-address=139.59.19.30:24224 \
		--log-opt tag={{.Name}} \
		-p 80:80 \
		-p 34567:34567 \
		--rm \
		--name injun \
		gw/injun:latest

show-images:
	docker images

show-containers:
	docker ps -a

delete-image:
	docker rmi gw/injun:latest

delete-container:
	docker rm injun_injunx_1