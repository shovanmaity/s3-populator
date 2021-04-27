# Introduction
In this example we will deploy a static website using nginx, minio and volume populator in kubernetes. This demo will be performed in local machine (minikube). 

# Steps
1. Start minikube with _AnyVolumeDataSource_ feature enable mode.
   ```bash
   minikube start --feature-gates=AnyVolumeDataSource=true
   ```
2. Install required components for volume populator.
   ```bash
   kubectl apply -f crd.yaml
   kubectl apply -f deploy.yaml
   ```
3. Install and configure any CSI driver. In this demo OpenEBS LVM CSI driver will be used.
   ```bash
   kubectl apply -f https://openebs.github.io/charts/lvm-operator.yaml
   ```
4. Configure a LVM volume group in local machine.
   ```bash
   truncate -s 10G /tmp/disk.img
   sudo losetup -f /tmp/disk.img --show
   sudo pvcreate /dev/loopX # replace X
   sudo vgcreate lvmvg /dev/loopX # replace X
   ```
5. Create a storage class for this LVM volume group.
   ```bash
   kubectl create -f - << EOF
   apiVersion: storage.k8s.io/v1
   kind: StorageClass
   metadata:
     name: openebs-lvmpv
   parameters:
     storage: lvm
     volgroup: lvmvg
   provisioner: local.csi.openebs.io
   EOF
   ```
6. Run minio as S3 source in local machine.
   ```bash
   docker run -p 9000:9000 \
   -e "MINIO_ROOT_USER=minioadmin" \
   -e "MINIO_ROOT_PASSWORD=minioadmin" \
   minio/minio server /data
   ```
7. Create a static file and push it to a minio bucket
   ```bash
   echo "Hello volume populator" > index.yaml
   ```
   ```bash
   AWS_ACCESS_KEY_ID=minioadmin \
   AWS_SECRET_ACCESS_KEY=minioadmin \
   AWS_REGION=us-east-1 \
   aws --endpoint-url http://192.168.0.190:9000 s3api create-bucket --bucket my-bucket
   ```
   ```bash
   AWS_ACCESS_KEY_ID=minioadmin \
   AWS_SECRET_ACCESS_KEY=minioadmin \
   AWS_REGION=us-east-1 \
   aws --endpoint-url http://192.168.0.190:9000 s3 cp . s3://my-bucket --recursive
   ```
8. Create a S3populator cr.
   ```bash
   kubectl create -f - << EOF
   apiVersion: example.io/v1
   kind: S3Populator
   metadata:
     name: s3-populator-1
     namespace: default
   spec:
     url: http://192.168.0.190:9000
     id: minioadmin
     secret: minioadmin
     region: us-east-1
     bucket: my-bucket
     key: /
   EOF
   ```
9. Deploy nginx with the data present in S3.
   ```bash
   kubectl create -f - << EOF
   apiVersion: apps/v1
   kind: StatefulSet
   metadata:
     name: web
   spec:
     serviceName: nginx
     replicas: 1
     selector:
       matchLabels:
         app: nginx
     template:
       metadata:
         labels:
           app: nginx
       spec:
         containers:
         - name: nginx
           image: k8s.gcr.io/nginx-slim:0.8
           ports:
           - containerPort: 80
             name: web
           volumeMounts:
           - name: www
             mountPath: /usr/share/nginx/html
     volumeClaimTemplates:
     - metadata:
         name: www
       spec:
         dataSource:
           apiGroup: example.io
           kind: S3Populator
           name: s3-populator-1
         storageClassName: openebs-lvmpv
         accessModes: [ "ReadWriteOnce" ]
         resources:
           requests:
             storage: 10Mi
   EOF
   ```
10. Create a service to access it
    ```bash
    kubectl create -f - << EOF
    apiVersion: v1
    kind: Service
    metadata:
      name: nginx
      labels:
        app: nginx
    spec:
      ports:
      - port: 80
        name: web
      clusterIP: None
      selector:
        app: nginx
    EOF
    ```
Now if we access the static website we can see that this is created using data present in S3.
NOTE - We can use any CSI driver for this volume populator.
