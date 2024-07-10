## Lamoi
Lamoi is a UI for ollama. It uses the Ollama API and stores conversations data in a local SQLite database

![lamoi](https://github.com/paganotoni/lamoi/assets/645522/de2541b3-de4a-4dfc-851c-23fc60067ab7)

### Getting started

#### Setup
Install The kit CLI.
```sh
go install github.com/leapkit/leapkit/kit@latest
```

#### Migrate database
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

### Building the application binary
To generate a distributable binary you can run the following command. This will build the TailwindCSS and the application.

#### Build the CSS
Building TailwindCSS with the tailo build command.
```sh
go run github.com/paganotoni/tailo/cmd/build@v1.0.7
```

#### Building the app
Building the application binary with Go.

```
go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app
```
