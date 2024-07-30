import json
import os
import time
from kubernetes import client, config, watch

def create_pod(obj_name, image):
    #print("-------------CREATE POD IS RUNNING-------------")
    config.load_kube_config()
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

def main():
    config.load_kube_config()
    v1 = client.CustomObjectsApi()

    resource_version = ''
    while True:
        stream = watch.Watch().stream(v1.list_namespaced_custom_object, group='example.com', version='v1', namespace='default', plural='mypods', resource_version=resource_version)
        for event in stream:
            obj = event['object']
            obj_name = obj['metadata']['name']
            spec = obj.get('spec', {})
            image = spec.get('image', 'busybox')  # Default to nginx:latest if no image specified
            resource_version = obj['metadata']['resourceVersion']
            if event['type'] == 'ADDED':
                create_pod(obj_name, image)

if __name__ == "__main__":
    main()
