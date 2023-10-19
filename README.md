# My Chatbot

This is a chatbot project written in Go that uses websockets to serve chats and includes HTML code for the frontend.

## Project Structure

The project has the following files:

- `public/index.html`: This file contains the HTML code for the chatbot's frontend.
- `public/script.js`: This file contains the JavaScript code for the chatbot's frontend.
- `main.go`: This file is the entry point of the application. It creates an instance of the chatbot and sets up the websocket server.
- `go.mod`: This file is the module definition file for the project. It lists the dependencies and the version of Go used.
- `go.sum`: This file contains the expected cryptographic checksums of the content of specific module versions.
- `README.md`: This file contains the documentation for the project.

## Usage

To run the chatbot, and run the following command:

```
go run main.go
```

This will start the server and the chatbot will be available at `http://localhost:8080`.

## Dependencies

The project uses the following dependencies:

- `github.com/gorilla/websocket`: A package for working with websockets in Go.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.