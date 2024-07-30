import kopf
import kubernetes.client
from kubernetes import config

config.load_kube_config()

executor_data_store = {}
name_task_update = ""


def update_task(taskId):
    global name_task_update
    api = kubernetes.client.CustomObjectsApi()
    body = {
        "spec": {
            "completed": 1,
            "taskId": taskId
        }
    }
    api.patch_namespaced_custom_object(
        group="example.com",
        version="v1",
        namespace="default",
        plural="taskupdates",
        name=name_task_update,
        body=body
    )
    print("Updated 'completed' spec of taskupdates CR.")
    print(f"name_task_update: {name_task_update}")

@kopf.on.create('example.com', 'v1', 'taskupdates')
def create_fn(spec, meta, **kwargs):
    print("~~~~~~~~~~~~~~~~~~~~ Reached operator - task update ~~~~~~~~~~~~~~~~~~~~~~~~")
    global name_tu
    completed = spec.get('completed')
    name_task_update = meta.get('name')
    print(f"Completed: {completed}")

@kopf.on.update('example.com', 'v1', 'taskupdates')
def update_fn(spec, new, **_):
    print("~~~~~~~~~~~~~~~~~~~~ Reached operator - task update modified ~~~~~~~~~~~~~~~~~~~~~~~~")
    completed = new.get('spec', {}).get('completed')
    print(f"Updated Completed: {completed}")


@kopf.on.create('example.com', 'v1', 'sparkexecmsgs')
def create_fn(spec, meta, **kwargs):
    print("~~~~~~~~~~~~~~~~~~~~ Reached operator - executor message ~~~~~~~~~~~~~~~~~~~~~~~~")
    type_of_call = spec.get('type')
    taskId = spec.get('taskId')
    executorId = spec.get('executorId')
    print(f"Type of call: {type_of_call}")
    print(f"Message: {spec.get('message')}")
    update_task(taskId)
    executor_data_store[executorId]["num_running_tasks"] -= 1
    print(f"Current Executor Data Store: {executor_data_store}")

@kopf.on.update('example.com', 'v1', 'sparkexecutors')
def update_fn(spec, old, new, meta, **kwargs):
    print("~~~~~~~~~~~~~~~~~~~~ Reached operator - executor updated ~~~~~~~~~~~~~~~~~~~~~~~~")
    old_num_running_tasks = old.get('spec', {}).get('numRunningTasks')
    new_num_running_tasks = new.get('spec', {}).get('numRunningTasks')
    print(f"Old numRunningTasks: {old_num_running_tasks}")
    print(f"New numRunningTasks: {new_num_running_tasks}")

    executor_id = spec.get('executorId')
    if executor_id:
        executor_data_store[executor_id]["num_running_tasks"] = new_num_running_tasks
        print(f"Executor ID: {executor_id} updated in store.")
        print(f"Data: {executor_data_store[executor_id]}")

    print(f"Current Executor Data Store: {executor_data_store}")

@kopf.on.create('example.com', 'v1', 'sparkexecutors')
def create_fn(spec, meta, **kwargs):
    print("~~~~~~~~~~~~~~~~~~~~ Reached operator - executor created ~~~~~~~~~~~~~~~~~~~~~~~~")
    type_of_call = spec.get('type')
    print(f"Type of call: {type_of_call}")

    executor = meta.get('name')
    executor_id = spec.get('executorId')
    executor_hostname = spec.get('executorHostname')
    os_info = spec.get('osInfo')
    java_version = spec.get('javaVersion')
    num_of_running_tasks = spec.get('numRunningTasks')
    print(f"Executor: {executor}")
    print(f"Executor ID: {executor_id}")
    print(f"Executor Hostname: {executor_hostname}")
    print(f"OS Info: {os_info}")
    print(f"Java Version: {java_version}")
    print(f"Num of running tasks: {num_of_running_tasks}")

    if executor_id:
        executor_data_store[executor_id] = {
            "hostname": executor_hostname,
            "os_info": os_info,
            "java_version": java_version,
            'num_running_tasks': num_of_running_tasks
        }
        print(f"Executor ID: {executor_id} added to store.")
        print(f"Data: {executor_data_store[executor_id]}")


    print(f"Current Executor Data Store: {executor_data_store}")

@kopf.on.create('example.com', 'v1', 'sparkrequests')
def create_fn(spec, meta, **kwargs):
    print("~~~~~~~~~~~~~~~~~~~~ Reached operator - executor request ~~~~~~~~~~~~~~~~~~~~~~~~")

    executorId = min(executor_data_store, key=lambda x: executor_data_store[x]['num_running_tasks'])
    print(f"Choosing this: Executor ID: {executorId}")
    executor_data_store[executorId]["num_running_tasks"] += 1

    api = kubernetes.client.CustomObjectsApi()
    body = {
        "spec": {
            "executorId": executorId,
        }
    }
    sparkrequest_name = meta.get('name')
    api.patch_namespaced_custom_object(
        group="example.com",
        version="v1",
        namespace="default",
        plural="sparkrequests",
        name=sparkrequest_name,
        body=body
    )
    print(f"Current Executor Data Store: {executor_data_store}")


if __name__ == '__main__':
    kopf.run(__file__)
