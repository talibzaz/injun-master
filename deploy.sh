#!/bin/bash

container_name=injun
semaphore_git_token=e8040157cee7a7cbb2cf764b81f8fd5195a2341b

# container checks
if docker ps -a --format '{{.Names}}' | grep -Eq "^${container_name}\$"; then
  if docker ps -aq -f status=running -f name=${container_name}; then
	docker stop ${container_name} && docker rm ${container_name}
  elif docker ps -aq -f status=exited -f name=${container_name}; then
    	docker rm ${container_name}
  fi
fi

# dep check
if ~/dep version; then
  echo running dep ensure...
  dep ensure
  echo DONE!
else
  echo go get -u github.com/golang/dep/cmd/dep
  go get -u github.com/golang/dep/cmd/dep
  echo running dep ensure...
  dep ensure
  echo DONE!
fi

make all