unexport VAGRANT_SERVER_URL

up:
	vagrant up

down:
	vagrant halt

ssh:
	vagrant ssh

client-test:
	go run cmd/main.go client ~/crate_test_dir
