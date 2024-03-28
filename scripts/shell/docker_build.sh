#!/usr/bin/env bash
echo "env-------------start"
echo "env-------------end"

cd ../../.. && docker buildx build --platform linux/amd64 \
  -f ./app/"$APP_RELATIVE_PATH"/Dockerfile \
  --build-arg CONFIG_FILE_NAME="$CONFIG_FILE_NAME" \
  --build-arg APP_RELATIVE_PATH="$APP_RELATIVE_PATH" \
  -t "$IMAGE" \
  .

if [ "$CI" = true ]; then
  echo "this image is push to remote server ..."
  docker push "$IMAGE"
else
  echo "this image is build locally and load into minikube ..."
  minikube image load "$IMAGE"
fi
