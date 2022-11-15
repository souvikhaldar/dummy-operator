# dummy-operator
This operator creates simple nginx pods corresponding to Spec and Updates status of Dummy object and logs certain provided message from the Spec.

# Test the operator

[Video Demonstration of this operator](https://youtu.be/QQzgsCWbe58)  

1. Clone this Repository.  
2. Run the operator using `make install run`.  
3. Create the namespace to test out the operator.`kubectl create ns dummy-test`   
4. Edit the sample config file by `vim config/samples/_v1alpha1_dummy.yaml`  
5. For eg. set the following value:  
```
apiVersion: souvikhaldar.in/v1alpha1
kind: Dummy
metadata:
  name: dummy-sample
spec:
  message: "I'm just a dummy"
  replicaCount: 6
  port: 80

```   
6. You can apply the config using: `kubectl apply -f config/samples/_v1alpha1_dummy.yaml -n dummy-test`  
7. Now you can check if required number of pods are running or not as specified in the spec object by `kubectl get pods -n dummy-test`.  


Alternatively, you can use the docker image https://hub.docker.com/r/souvikhaldar11/dummy-operator



