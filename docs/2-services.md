# Tutorial : first steps with OTO (part 2)

Hello and welcome to the first steps tutorial of OTO go service.
Lets learn how to download and use OTO.


# Introduction

Currently, OTO is under heavy development and don't have a proper tagged version on the main branch. Which mean you can't call the go package easily go mod tidy or go get.

In order to test OTO before the first official release, I suggest you to download directly the repository and test OTO in the given main.go file at the root of the project.

So the first step to use OTO is to git clone the project :
```bash
git clone https://github.com/Bl4omArchie/oto
```

One it is done, go in the **oto/** folder.


# Services

Now, lets begin the initialization of the required services to make OTO work

OTO is running with the following services :
- temporal server
- postgres database and pgadmin
- rest API
- web dashboard

In order to setup everything rigth, you can call the docker compose file.
Follow those steps :
1. Install Docker
2. Create a .env file with the following variables :
```
POSTGRES_HOST=postgres
POSTGRES_DB=your_database
POSTGRES_USERyour_user
POSTGRES_PASSWORD=your_password
POSTGRES_PORT=5432
POSTGRES_SEED=postgres

TEMPORAL_HOST=:7233
TEMPORAL_NAMESPACE=default
```

> Take care of not interfering with existing services on your machine because the port number that can be the same. 

3. Launch the following command :
```bash
docker compose up -d
```

If everything turn green and says **started**, then your good.


# Web dashboard

You can now click on the following address to access the different web dashboard :
- OTO dashboard -> http://localhost:2929
- Pgadmin -> http://localhost:5050
- Temporal web ui -> http://localhost:8080

The API is also available here : http://localhost:1515

Next part of the guide [here](3-code.md)
