# Forum
This project consists in creating a web forum that allows :

crate users, create and tag posts, comment posts, react to posts and comments, communicate by comments, filter posts
## Usage

`make run` <br> -> build and run docker image on port :8000 <br> or <br>
`docker image build -f Dockerfile -t forum-image .` <br> `docker run -p 8000:8000 -d --name forum-container forum` <br>
-> build and run container on port :8000 <br>
`make del`<br> -> stop container, delete it anf all images
## Contributors
Sserbulo<br />
Dantonen