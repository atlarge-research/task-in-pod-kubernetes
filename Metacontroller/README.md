# Metacontroller framework #

## Sequence of commands to execute: ##

<sub>(Assuming the minikube is running already)</sub>

### 1. Run the sync.py or sync.go code: ###

- sync.py: ```python sync.py```

- sync.go: ```go run "your-path\sync.go" -kubeconfig="$env:USERPROFILE\.kube\config```

### 2. In another terminal and run the following commands: ###
1. ```kubectl apply -f crd.yaml```
2. ```kubectl apply -f controller.yaml```
3. ```kubectl apply -f obj.yaml```

<sub>If running sync.go - replace **```command: ["python", "sync.py"]```** with **```command: ["./sync"]```** in controller.yaml file.</sub>

### 3. Verify and clean: ###
 - To check the pods, run ```kubectl get pods```

 - To delete, run ```kubectl delete -f obj.yaml```

 - To delete all pods, run ```kubectl delete pods --all```
