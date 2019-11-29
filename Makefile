unexport VAGRANT_SERVER_URL

up:
	vagrant up

down:
	vagrant halt

ssh:
	vagrant ssh

client-test:
	go run cmd/client/main.go ~/crate_test_dir
