# VitaCare


## How to run the project
1. Clone the project
2. Install [devbox](https://www.jetify.com/docs/devbox/installing_devbox/)
3. Be sure you have installed Makefile
4. Install docker and optional have docker desktop

## Set up the project
1. Have a instance of postgres running on your machine, you can do it with docker with the next command 
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres \
-p 5432:5432 -d postgres
```
2. Copy the `config/config.example.yml` file and rename it to `config/config.yml` and modify the values to match your local configuration


3. Go to the project folder and run the next command
```bash
    make shell && make db-up 
```
this command will init the devbox shell and run all the migrations for the project

4. Finally run the project with the next command
```bash
    make dev
```

## How to run the tests
Go to the project folder and run the next command
```bash
    make test
```

### helpful commands
to know all the commands available run the next command
```bash
    make help
```
