# Forum
This project consists in creating a web forum that allows : <br>

create and login users, save sesions, create and tag posts, comment posts, react to posts and comments, communicate by comments, filter posts
## Usage

`make run` <br> -> build and run docker image on port :8000 <br> or <br>
`docker image build -f Dockerfile -t forum-image .` <br> `docker run -p 8000:8000 -d --name forum-container forum` <br>
-> build and run container on port :8000 <br>
`make del`<br> -> stop container, delete it anf all images
## Contributors
[@sunf1ower113](https://github.com/Sunf1ower113)<br />

## Notes
This is a graduation project from Alem School in the direction of backend developer for Go, completed on April 1, 2023. <br>
The project was added in the form and with the errors that were in the code review to show progress since graduation.<br>
