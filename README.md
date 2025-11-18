# OTO - workflow automation service

OTO is a golang service for workflow automation.


Currently, OTO is providing a single but powerful feature : **RunJob**. 
This task allow you to run binaries automatically with pre defined values.

# How to use OTO ?

OTO is a go service that comes up with a sqlite database, a go package, a web interface and a restfulAPI. 
Which give you three ways for using OTO :
- Go package : call the public functions yourself in a go project.
- Restful API : the service provides a resful api powered by Gin so you can use OTO through command line.
- Web interface : finally, OTO provide a web interface connected to the API so you can manage your workflows more easily.
