#!/bin/sh

#docker --version

docker build -t go-doc:0.0.2 -f /home/wojtek/workspace/WeatherHealthChecker/Dockerfile.multistage /home/wojtek/workspace/WeatherHealthChecker
docker tag go-doc:0.0.2 <aws account>.dkr.ecr.eu-west-1.amazonaws.com/go-doc:0.0.2
docker push <aws account>.dkr.ecr.eu-west-1.amazonaws.com/go-doc:0.0.2
kubectl rollout  restart deployment/goapp17

#kubectl apply -f goapp-deployment.yaml




