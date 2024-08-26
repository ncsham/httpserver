# httpserver
Basic HTTP Server in Golang

# Port the Image to Minikube
- https://minikube.sigs.k8s.io/docs/handbook/pushing/#1-pushing-directly-to-the-in-cluster-docker-daemon-docker-env
- `eval $(minikube docker-env)`
- `docker build -t httpserver .`
- Apply the k8s manifest (create staging namespace if doesn't exist or modify the manifest)
  - `kubectl apply -f k8s.yaml`