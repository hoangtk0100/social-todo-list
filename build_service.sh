IMAGE_NAME=social-todo-service
CACHED_BUILD=$1

if [[ -n "$CACHED_BUILD" ]]; then
    echo "Docker building cache image..."
    docker rmi ${IMAGE_NAME}-cache ${IMAGE_NAME}
    docker build -t ${IMAGE_NAME}-cache -f Dockerfile-cache .
fi

echo "Docker building main image..."
docker build -t ${IMAGE_NAME}:latest -f Dockerfile-with-cache .

echo "Done!"
