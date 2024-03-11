build:
	docker build -f Dockerfile -t forum_image .

run:	build
		docker image ls
		docker run -d -p 8000:8000 --name forum_container forum_image
		docker ps -a
		docker logs --follow forum_container

del image:
	docker stop forum_container
	docker rm forum_container
	docker image rm forum_image
	docker image prune