# Drone CI K8s manifest

The following repo holds the k8s manifest to deploy our CI instance

Keep in mind that in order to expose this service propely, you'll need to create an ingress entry from the gke dashboard.
Also, the GH app secret needs to be defined:

```
kubectl create secret generic droneci --from-literal=github-client-secret='YourSecret'
```
