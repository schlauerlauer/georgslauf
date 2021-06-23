IMAGE=georgslauf
PORT=3000

build:
	podman build -t $(IMAGE) .
run:
	podman run -it --rm -p $(PORT):$(PORT) $(IMAGE)
