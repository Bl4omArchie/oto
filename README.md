# OTO - workflow automation service

OTO is a golang service for workflow automation.


Currently, OTO is providing a single but powerful feature : **RunJob**. 
This task allow you to run binaries automatically with pre defined values.

# How to use OTO ?

There is several services in OTO all of them launched from the docker compose
- Restful API : powered by Gin
- Web dashboard : simple html/css/js dashboard for easy management
- Go package : get OTO package to use the service directly in your own code
- Postgresql DB : for data persistency (linked to Temporal)
- Temporal server : launch workers and run workflows
- Temporal web UI : monitoring of workflows

To launch everything use the following command :
```sh
docker compose up -d
```
Before it you must  set a .env file with credentiels for the postgresql DB.


# FME : flag matching engine

OTO has recently integrated the `lvlath` library in order to create a flag matching engine. Let me do a quick introduction of this new piece of code.

**What is the flag matching engine ?**

Basically, a **Command** is a set of **Parameter**. But not every parameters can match each other into one Command.
For instance if a command with parameter **b** won't work with parameter **a**, there is no interest in letting the user do so.

OTO is now able to efficiently refuse such commands with the flag matching engine, a program that verify dependencies and conflict between parameters.

Here si the full notebook for further details and in coming fearures [click here](https://github.com/Bl4omArchie/flag_matching)

# Integration of Temporal

The very core of my service is now done and I can now think about scheduling and monitoring larger configuration of binaries with a more large amount of commands.
To do so, I have to handle concurrency, events log and more to make a robust solution.
As a solution, I've started the integration of Temporal.io in the `dev` branch where I launch, through a docker compose multiple services like temporal.io, temporal ui, postgresl and my API.

Why Temporal ?
- Event logging
- Metrics
- Workflows and workers for my jobs
- Mature and globally used framework

# Roadmap

Current work :
- [] import/export config to json
- [] docker compose deployement
- [] 