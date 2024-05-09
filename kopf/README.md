# Kopf Framework #

## Sequence of commands to execute: ##

<sub>(Assuming the minikube is running already)</sub>

### 1. Run the kopf.py file: ###

- ```kopf run kopf.py --verbose```


### 2. In another terminal and run the following commands: ###
1. ```kubectl apply -f crd.yaml```
3. ```kubectl apply -f obj.yaml```


### 3. Verify and clean: ###
 - To check the pods, run ```kubectl get pods```

 - To delete, run ```kubectl delete -f obj.yaml```

 - To delete all pods, run ```kubectl delete pods --all```
