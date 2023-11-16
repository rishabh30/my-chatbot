<img width="442" alt="image" src="https://github.com/rishabh30/my-chatbot/assets/39810768/81d90812-1701-4297-8ea5-6c15ce908338">
# My Chatbot

This is a prototype for a chatbot written in Go that allows users to save and retrieve images. Also, is enables the chatbot to discuss with the user using off-the-shelf language model solutions. 

## Project Structure

The project has the following files:

- `public/index.html`: This file contains the HTML code for the chatbot's frontend.
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
