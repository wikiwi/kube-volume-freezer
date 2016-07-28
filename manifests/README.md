# Deploy kube-volume-freezer
This folder contains an example deployment of kube-volume-freezer with token protected Minions and API server.

# Walkthrough
- Set your base64 encoded tokens in the `kcf-secret.yaml` file.
- Add `kvf-secret.yaml`, `kvf-daemonset.yaml`, `kvf-deployment.yaml` and `kvf-svc.yaml` to your Kubernetes.
