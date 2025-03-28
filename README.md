
## Getting Started

### For Istio, make sure once GCP infra is up, apply istio custom profile in tf-gcp-project repo. 

```sh
- gcloud container clusters get-credentials dev-gke-cluster --region us-east1

- istioctl install -f custom-istio.yaml —skip-confirmation

- cd helm/api-server
- helm install webapp .

```

### For cert manager setup on GKE cluster: 

```sh
helm repo add jetstack https://charts.jetstack.io
helm repo update

# Install cert-manager
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.17.0 \
  --set installCRDs=true
```

### Create Helm release

``` sh
- cd helm/api-server
- helm install webapp .

```

### Update Istio's Ingress Gateway's External IP as an A record in Cloud DNS GCP console

### Test out webapp endpoints and functionaly using Postman APIs with HTTPS protocol

### Decrypt using SOPS: 
```sh

sops -d secrets-enc.yaml > secrets-denc.yaml
```

