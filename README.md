## Newspeak - Project Info

### About
This is a still a work in progress, it's primarily a project to tech myself some new  techniques.

The goal is to create a web based chat app in which users can create and join rooms based on geographical location. The backend is written in Golang and it's designed around WebSockets and with concurrency support. The frontend is being done in VueJS.

### Docker info:

Building the image: `docker build -t newspeak-app .`

Running the image: `docker run -d -p 8000:8000 newspeak-app`