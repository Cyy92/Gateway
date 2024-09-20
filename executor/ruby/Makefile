registry=10.0.1.150:5000/cyy
version=1.0
ruby_version=2.6.5

all: docker-build docker-push

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-ruby --build-arg RUBY_VERSION=$(ruby_version) .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-ruby
