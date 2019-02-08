.PHONY: build clean deploy

build:
	dep ensure
	docker build -t ygrene/iam-eks-user-mapper .
	DOCKER_CONTENT_TRUST=1 docker push ygrene/iam-eks-user-mapper:latest

deploy: build
	kubectl apply -f kubernetes/
