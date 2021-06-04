# synk8s

Synchronize k8s native resources between the source and destination cluster. Deletion of the source instance will trigger deletion of the destination instance as well.

## Binary & Run Binary

```bash
# Build the binary. Here the default values are also given.
# Supported group-version can be found in this README.
# Modify these four variables to suit your repository.
# For the most resources, `FIELD` will be `Spec`.
make \
GROUPVERSION=corev1 \
KIND=Secret \
FIELD=Data

# Use the above defaults
make

# Run the binary.
# Assign a namespace of interest. Make sure the valid paths to the source and destination clusters are given.
bin/synk8s --namespace cloud-gateway-server --source /tmp/source --dest /tmp/dest
```

## Build & Run Docker Image

```bash
# Build the image. As above, the default values are also given.
make docker-build \
GROUPVERSION=corev1 \
KIND=Secret \
FIELD=Data

# Similarly, use the above defaults
make docker-build

# Run the image. As above, make sure `/tmp/source` and `/tmp/dest` are valid paths and assign a namespace of interest.
docker run -v /tmp/source:/source:ro -v /tmp/dest:/dest:ro --network host ghcr.io/metal-stack/synk8s:latest --namespace=cloud-gateway-server
```

## Supported group-version

```bash
admissionregistrationv1
admissionregistrationv1beta1
internalv1alpha1
appsv1
appsv1beta1
appsv1beta2
authenticationv1
authenticationv1beta1
authorizationv1
authorizationv1beta1
autoscalingv1
autoscalingv2beta1
autoscalingv2beta2
batchv1
batchv1beta1
batchv2alpha1
certificatesv1
certificatesv1beta1
coordinationv1beta1
coordinationv1
corev1
discoveryv1alpha1
discoveryv1beta1
eventsv1
eventsv1beta1
extensionsv1beta1
flowcontrolv1alpha1
flowcontrolv1beta1
networkingv1
networkingv1beta1
nodev1
nodev1alpha1
nodev1beta1
policyv1beta1
rbacv1
rbacv1beta1
rbacv1alpha1
schedulingv1alpha1
schedulingv1beta1
schedulingv1
storagev1beta1
storagev1
storagev1alpha1
```
