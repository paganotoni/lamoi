## Lamoi

Lamoi is a UI for ollama. It uses the Ollama API and stores conversations data in a local SQLite database

### Getting started

### Setup
Install The kit CLI.
```sh
go install github.com/leapkit/leapkit/kit@latest
```

### Migrate database
Kit takers care of this part:
```
kid db migrate
```

### Running the application
To run the application in development mode execute:
```sh
kit serve
```

And open `http://localhost:3000` in your browser.

### Building the application

To generate a distributable binary you can run the following command. This will build the TailwindCSS and the application.

```sh
# Building TailwindCSS with tailo
> go run github.com/paganotoni/tailo/cmd/build@a4899cd

# Building the app
> go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app
```
