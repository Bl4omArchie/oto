# Changelog


**02/12/25** :
- Add `envPath` parameter for NewInstanceOto() : you can now specify DB you want to use
- ImportParameters() can import a json file of parameters into the DB
- New struct for json data : ParameterRaw

**25/11/25** :
- Add dockerfile for API
- Replace sqlite with postgres
- Add docker compose for postgres, pgadmin, temporal, temporal-ui and the api server
- Add test package
- Add .env file for a more secure database connection

**23/11/25** :
- Integration of (FME) Flag Matching Engine for fast dependencies and conflicts check for Commands
- Typo : 
    - Executable becomes Binary
    - JobCommand becomes Job
- New functions to fetch more easily models

**12/11/25** : 
- Add restfulAPI
- Add web dashboard squeleton
- Add cmd/ folder for starting API and dashboard
- Work on : flag matching mecanism

**11/11/25** :
- First activity done : RunJobCommand
- Database schema done and tested
- Add documentation
