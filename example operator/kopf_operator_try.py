import kopf

@kopf.on.startup()
def on_startup_fn(**kwargs):
    print("Hello, world! Operator started.")

@kopf.on.create('kopfexamples')
def create_fn(spec,**kwargs):
    print("-----CREATED-----  "+str(spec))

if __name__ == "__main__":
    kopf.run(on_startup_fn)
    kopf.run(create_fn)