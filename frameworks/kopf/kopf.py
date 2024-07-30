import kopf
from kubernetes import client, config


@kopf.on.create('example.com', 'v1', 'mypods')
def create_pod_fn(spec, **kwargs):
    #print("-------------CREATE POD IS RUNNING-------------")
    obj_name = kwargs['body']['metadata']['name']
    image = spec.get('image', 'nginx:latest')  # Default to nginx:latest if no image specified

    config.load_kube_config()  # Load kubeconfig for local testing, use config.load_incluster_config() for in-cluster
    api = client.CoreV1Api()

    pod_manifest = {
        'apiVersion': 'v1',
        'kind': 'Pod',
        'metadata': {'name': f'pod-{obj_name}'},
        'spec': {
            'containers': [
                {
                    'name': 'my-container',
                    'image': image,
                    'ports': [{'containerPort': 80}]
                }
            ]
        }
    }

    api.create_namespaced_pod(namespace='default', body=pod_manifest)
