registry=10.0.1.150:5000/cyy
version=1.0
java_version=8u242

all: docker-build docker-push

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-java --build-arg JAVA_VERSION=$(java_version) .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-java
