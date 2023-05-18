IMAGE=georgslauf
PORT=8080

build:
	podman build -t $(IMAGE) .
run:
	podman run -it --rm -p $(PORT):$(PORT) $(IMAGE)
