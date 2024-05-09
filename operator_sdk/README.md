# Operator SDK framework #
Here, there are two operators in the folder - the one that creates a pod and the one that simply prints ("Hello, I am an operator"). 
The only differences in the process between the two are:
- filenames
- code inside *.go* file

Hence, replace the names of the files accordingly.

### Steps for creating and running an operator: ###
1. Create directory: ```mkdir pod_creator_operator```
2. Go to that dir: ```cd pod_creator_operator``` 
3. Initialize new operator project: ```operator-sdk init --domain=your.domain --repo=github.com/yourusername/pod-creator-operator```
4. Create a new API for Pod resource: ```operator-sdk create api --group=example --version=v1 --kind=Pod```
5. Modify file */pod_creator_operator/internal/controller/pod_controller.go* - required code is in *sdk_pod_controller.go* file
6. Build and deploy an operator:
   - make generate
   - make manifests
   - make install
   - make docker-build docker-push *(You might need to login to your docker account with ```docker login```, then run ```docker tag example.com/pod-creator-operator:0.0.1 username/pod-creator-operator:0.0.1```and ```docker push username/pod-creator-operator:0.0.1```)*
   - make deploy
7. Create a new yaml file (the one like sdk_pod.yaml)
8. Run ```kubectl apply -f sdk_pod.yaml```
9. To see the pods, run ```kubectl get pods```
10. To clean, run ```kubectl delete -f sdk_pod.yaml```
