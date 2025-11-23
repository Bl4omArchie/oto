# OTO - workflow automation service

OTO is a golang service for workflow automation.


Currently, OTO is providing a single but powerful feature : **RunJob**. 
This task allow you to run binaries automatically with pre defined values.

# How to use OTO ?

OTO is a go service that comes up with a sqlite database, a go package, a web interface and a restfulAPI. 
Which give you three ways for using OTO :
- Go package : call the public functions yourself in a go project.
- Restful API : the service provides a resful api powered by Gin so you can use OTO through command line.
- Web interface : finally, OTO provide a web interface connected to the API so you can manage your workflows through the dashboard.

# FME : flag matching engine

OTO has recently integrated the `lvlath` library in order to create a flag matching engine. Let me do a quick introduction of this new piece of code.

**What is the flag matching engine ?**

Basically, a **Command** is a set of **Parameter**. But not every parameters can match each other into one Command.
For instance if a command with parameter **b** won't work with parameter **a**, there is no interest in letting the user do so.

OTO is now able to efficiently refuse such commands with the flag matching engine, a program that verify dependencies and conflict between parameters.

Here si the full notebook for further details and in coming fearures [click here](https://github.com/Bl4omArchie/flag_matching)
