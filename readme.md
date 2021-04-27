# S3-Populator
S3 populator is a kubernetes volume populator which helps us to create a volume with data present in S3 bucket. It can pull data present in a bucket or sub directory of a S3 bucket. Please make sure _AnyVolumeDataSource_ feature is enabled in kubernetes cluster. You can do it in minikube using ->

_minikube start --feature-gates=AnyVolumeDataSource=true_ .

To install other components run below commands
```bash
kubectl apply -f crd.yaml
kubectl apply -f deploy.yaml
```

Here is an example of S3 populator cr
```yaml
apiVersion: example.io/v1
kind: S3Populator
metadata:
  name: s3-populator-1
  namespace: default
spec:
  url: http://192.168.0.190:9000 # <- S3 server url
  id: minioadmin # <- id to access the bucket
  secret: minioadmin # <- secret/key to access the bucket
  region: us-east-1 # <- region for the bucket
  bucket: b-001 # <- bucket name
  key: / # <- path that we want to download from a bucket.
  # If we want to download full bucket we can provide /.
  # If we want to download a sub dir from a bucket then we
  # can provide the sub dir name ie - /dir-1 .

```

## Example
Here are some example how we can use S3 populator
1. [Static website.](https://github.com/shovanmaity/s3-populator/tree/master/example/static-website.md)
