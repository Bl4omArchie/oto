# Changelog

**15/12/25** :
- Correction into README
- Add comments above functions in oto.go for clarity
- Add new feature called Routine. Your now able to decide of the order of execution of your jobs.
- Done with the temporal workflow : RunRoutine
- Changed the initial workflow name RunJob to RunRoutine

**10/12/25** :
- Integration of Atlas for automatic database migration
- DB update : [migrations/20251210120506.sql](migrations/20251210120506.sql)
- Typo fix in oto.go (bin -> exec)

**09/12/25** :
- Pre-Release + demo with openSSL
- User manual in docs/
- Update docker compose

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
    - Executable becomes Executable
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
