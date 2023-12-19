# FORUM

- [Audit instructions](https://01.kood.tech/git/adubtsov/forum#audit-instructions)
- [Task description](https://github.com/01-edu/public/blob/master/subjects/forum/README.md)
- [Audit on 01-edu](https://github.com/01-edu/public/blob/master/subjects/forum/audit/README.md)

### Collaborators

- [Alexey Dubtsov](https://01.kood.tech/git/adubtsov)
- [Ivar Mihhailov](https://01.kood.tech/git/imihhail)
- [Kaarup Vares](https://01.kood.tech/git/kvares)
- [Tanel Soidla](https://01.kood.tech/git/TanelS)

## Audit instructions:

### Alternative 1:

1. Open terminal. Clone this repo using following command:

```console
git clone https://01.kood.tech/git/adubtsov/forum
```

2. Change current folder to folder with the project:

```console
cd ./forum/
```

3. Run the server with command:

```console
go run .
```

4. Open this page in web browser:

```console
http://localhost:8080/home
```

5. Check our project using audit questions:

- [Link](https://github.com/01-edu/public/blob/master/subjects/forum/audit/README.md)

### Alternative 2:

This audit is done with using Docker.

#### What is Docker

Docker is like a virtual box for your code. Imagine you have a piece of software, and it needs a specific environment to run smoothly (certain libraries, system settings, etc.). Docker helps you package your software along with all its needs into a 'container'. This way, wherever you move this container, be it your laptop, a friend's computer, or a server, the software will run in the same way, because it has everything it needs inside the container.

So, in simple terms, Docker is a tool that helps you create and use these containers easily, ensuring that your software runs the same, no matter where it is.

#### Technologies

- **Language:** Go
- **Containerization:** Docker

#### Prerequisites

Docker Desktop should be installed on your machine.

Instructions how to do that:
- [Linux Mint, Ubuntu](https://docs.docker.com/desktop/install/ubuntu/)
- [MacOS](https://docs.docker.com/desktop/install/mac-install/)
- [Windows](https://docs.docker.com/desktop/install/windows-install/)

[General instructions](https://docs.docker.com/desktop/)



#### Audit instructions:

1. Open terminal, clone this repo using command

```console
git clone https://01.kood.tech/git/adubtsov/forum
```

2. Change current folder to folder with the project:

```console
cd ./forum/
```

3. Build Docker Image:

```console
docker image build -t my-dockerize-app .
```

It will take some time because of downloading needed packages.

4. After building the image, list all Docker images:

```console
docker images
```

5. Run Docker Container:

```console
docker container run -p 8080:8080 my-dockerize-app
```

After running this command, you should be able to access the web application by navigating to the link in your web browser:
[localhost:8080/home](http://localhost:8080/home)

6. Check our project using audit questions:

- [Link](https://github.com/01-edu/public/blob/master/subjects/forum/audit/README.md)

```console
https://github.com/01-edu/public/blob/master/subjects/forum/audit/README.md
```

7. List all running containers:

```console
docker ps -a
```

8. Stop running container:

```console
docker stop <container ID>
```

9. If you want to remove the stopped container, you can use the docker rm command:

```console
docker rm <container ID>
```

10. If you want to stop and remove the container in a single command, you can use the -f option with docker rm:

```console
docker rm -f <container ID>
```

11. With this command you can delete all containers located on your PC

Be careful!

```console
sudo docker system prune -a
```

#### Gratitude
Thank you!