# Task-aware scheduler #
The folder *spark-source-code* contains the modified version of the Spark source code. The main changes are in _core/src/main/scala/org/apache/spark/executor/Executor.scala_
and _core/src/main/scala/org/apache/spark/scheduler/cluster/CoarseGrainedSchedulerBackend.scala_. 

Note that the general _pom.xml_ file is also changed to include the kubernetes client library.

The code of the scheduler is inside the _spark-operator_ folder.

## Create a distribution of modified Spark ##
Steps to create a distribution of Spark:
<ol>
  <li> Run: <code>./dev/make-distribution.sh --name custom-spark-1 --tgz -DskipTests -Pkubernetes</code> </li>
  <li> Extract the tar: <code>sudo tar -xzf spark-3.5.1-bin-custom-spark-1.tgz -C /opt</code></li>
</ol>

## Run an application locally with Task-aware scheduler ##
This is the workflow to run an application locally but with Kubernetes support. We use Spark Pi as an example application.
<ol>
  <li> Create docker image of custom Spark distribution: in <code>/opt/your-custom-spark</code> run <code>docker build --network host -t image-name -f kubernetes/dockerfiles/spark/Dockerfile .</code></li>
  <li> Start Minikube (or similar) cluster: <code>minikube start --kubernetes-version=v1.25.3</code>. We use the Kubernetes version 1.25.3 because it is compatible with Spark and the Kubernetes client version used in the source code.</li>
  <li> Load the docker image to Kubernetes cluster: <code>minikube image load simple-app:latest</code></li>
  <li> In *spark-operator* folder run: <code>kubectl create serviceaccount spark</code></li>
  <li> Apply rbac: <code>kubectl apply -f rbac2.yaml</code></li>
  <li> Apply cluster role binding: <code>kubectl apply -f cluster-role-binding.yaml</code></li>
  <li> Apply CRDs: </li>
  <ul> <code>kubectl apply -f crd_executor.yaml</code></ul>
  <ul> <code>kubectl apply -f crd_executor_msg.yaml</code></ul>
  <ul> <code>kubectl apply -f crd_task_update.yaml</code></ul>
  <ul> <code>kubectl apply -f crd_spark_request.yaml</code></ul>
  <li> In separate tab, run: <code>kopf run sync4.py</code></li>
  <li> Run application <code>bin/spark-submit --master k8s://https://yourAddress --deploy-mode cluster --name simple-app --class org.apache.spark.examples.SparkPi --conf spark.executor.instances=2 --conf spark.kubernetes.container.image=yourImage:latest --conf spark.kubernetes.container.image.pullPolicy=IfNotPresent --conf spark.kubernetes.authenticate.driver.serviceAccountName=spark local:///opt/spark/examples/jars/spark-examples_2.12-3.5.1.jar</code></li>
  <li> After the application has finished, terminate the scheduler.</li>
</ol>

## Run TPC-DS benchmarking on a cluster using Sharebench framework ##
<ol>
  <li> Configure the cluster and start Kubernetes</li>
  <li> Create kube configuration: <code>scp -i $HOME/.ssh/id_rsa_continuum cloud_controller_$USER@yourIPAddress:/home/cloud_controller_$USER/.kube/config ~/.kube/ && kubectl auth can-i create deployments</code>
  <li> Create Spark service account: <code>kubectl create serviceaccount spark && \
kubectl create clusterrolebinding spark-role --clusterrole=edit --serviceaccount=default:spark --namespace=default</code></li>
  <li> Copy the distribution to the custom spark to the sharebench/spark folder </li>
  <li> Modify the Dockerfile inside the custom Spark folder to the one </li>
  <li> Create the docker image: <code>docker build --network host -t image-name -f kubernetes/dockerfiles/spark/Dockerfile .</code></li>
  <li> Modify the Dockerfile in the Sharebench framework like in the file </li>
  <li> Build and push docker image with <code>python3 scripts/image.py</code></li>
  <li> Generate the data: <code>spark/your-custom-spark/bin/spark-submit --class "ShareBench" --properties-file ./spark-defaults.conf --deploy-mode cluster local:///opt/sharebench/sharebench_2.12-1.0.jar datagen s3a://data /opt/sharebench/tpcds-bin/</code></li>
  <li> Generate metadata: <code>spark/your-custom-spark/bin/spark-submit --class "ShareBench" --properties-file ./spark-defaults.conf --deploy-mode cluster local:///opt/sharebench/sharebench_2.12-1.0.jar metagen s3a://data</code></li>
  <li> In a separate tab, run: <code>kopf run sync4.py</code></li>
  <li> Run benchmark: <code>spark/your-custom-spark/bin/spark-submit --class "ShareBench" --properties-file ./spark-defaults.conf --deploy-mode cluster local:///opt/sharebench/sharebench_2.12-1.0.jar queries_tpcds 3</code></li>
  <li> To get the logs of the most recent driver: <code>kubectl logs $(kubectl get pods -A --sort-by=.metadata.creationTimestamp | grep driver | tail -n 1 | awk '{print $2}') -n $(kubectl get pods -A --sort-by=.metadata.creationTimestamp | grep driver | tail -n 1 | awk '{print $1}') --follow</code></li>
  <li> To get the most recent driver outputs: <code>kubectl logs $(kubectl get pods -A --sort-by=.metadata.creationTimestamp | grep driver | tail -n 1 | awk '{print $2}') -n $(kubectl get pods -A --sort-by=.metadata.creationTimestamp | grep driver | tail -n 1 | awk '{print $1}') --follow | grep -v "^[0-9]\{2\}/[0-9]\{2\}/[0-9]\{2\} [0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\}"</code></li>
</ol>
