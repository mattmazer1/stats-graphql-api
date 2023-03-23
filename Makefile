.PHONY: start

start:
	gnome-terminal -- bash -c "nats-server; exec bash" & gnome-terminal -- bash -c "go run server.go; exec bash"
