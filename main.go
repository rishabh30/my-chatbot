package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"my_chatbot/server"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiURL     = "https://api.openai.com/v1/chat/completions"
	authBearer = "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"
)

type RequestPayload struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type UserMessage struct {
	Userid  int    `json:"userid"`
	Message string `json:"message"`
}

type ResponsePayload struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type Message struct {
	UserID      int    `json:"userid"`
	UserMessage string `json:"usermessage"`
	Response    string `json:"response"`
	Timestamp   string `json:"timestamp"`
	Path        string `json:"path"`
}

var db *sql.DB

func initializeDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./messages.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY, 
    userid INTEGER, 
    usermessage TEXT, 
    response TEXT, 
    timestamp TEXT,
	path TEXT
)`)
	if err != nil {
		panic(err)
	}
}

func storeMessage(userid int, usermessage, response string, timestamp time.Time, path string) {
	_, err := db.Exec("INSERT INTO messages (userid, usermessage, response, timestamp, path) VALUES (?, ?, ?, ?, ?)", userid, usermessage, response, timestamp.Format("2006-01-02 15:04:05"), path)
	if err != nil {
		panic(err)
	}
}

func replyToUser(w http.ResponseWriter, r *http.Request) {
	var userMsg UserMessage

	err := json.NewDecoder(r.Body).Decode(&userMsg)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if message is from the particular user
	if userMsg.Userid == 1 {
		prompt := userMsg.Message

		responseText := callOpenAI(prompt)

		storeMessage(userMsg.Userid, prompt, responseText, time.Now(), "")
		w.Write([]byte(responseText))
		w.WriteHeader(http.StatusOK)
	}
}

func callOpenAI(prompt string) string {
	data := RequestPayload{
		Model:     "gpt-3.5-turbo",
		Prompt:    prompt,
		MaxTokens: 150,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "Error occurred"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+authBearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error occurred"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response ResponsePayload
	json.Unmarshal(body, &response)

	if len(response.Choices) > 0 {
		return response.Choices[0].Text
	}
	return "No reply generated"
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Limit the data size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	msg := r.FormValue("message")
	userid := r.FormValue("userid")
	file, header, err := r.FormFile("file")
	// filess := r.MultipartForm.File["file"]
	// for _, filee := range filess {
	// 	fmt.Println(filee.Filename)
	// }

	if (err != nil && msg == "") || userid == "" {
		http.Error(w, "Both message and file are missing or userid is missing", http.StatusBadRequest)
		return
	}

	responseMsg := ""

	if file != nil {
		defer file.Close()

		dstPath := path.Join("files", header.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Failed to create the file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to save the file", http.StatusInternalServerError)
			return
		}

		responseMsg += fmt.Sprintf("Successfully uploaded the file: %v\n", header.Filename)
	}

	if msg != "" {
		responseMsg += fmt.Sprintf("Received message: %v\n", msg)
	}

	responseText := callOpenAI(msg)
	id, err := strconv.Atoi(userid)
	if err != nil {
		http.Error(w, "Please pass the user id in correct format", http.StatusBadRequest)
		return
	}
	storeMessage(id, msg, responseText, time.Now(), "")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseText))
}

func getMessage(userid int) []Message {
	rows, err := db.Query("SELECT userid, usermessage, response, path, timestamp FROM messages WHERE userid = ? order by timestamp desc limit 3", userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.UserID, &msg.UserMessage, &msg.Response, &msg.Path, &msg.Timestamp)
		if err != nil {
			panic(err)
		}
		messages = append(messages, msg)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return messages
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("userid")
	if userIDParam == "" {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}
	messages := getMessage(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	initializeDatabase()
	// http.HandleFunc("/send_message", replyToUser)
	http.HandleFunc("/send_message", uploadHandler)
	http.HandleFunc("/get_message", getMessageHandler)
	http.HandleFunc("/ws", server.WebsocketHandleConnection)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	// Start the server
	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// func callOpenAI(prompt string) string {
// 	client := openai.NewClient("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")

// 	resp, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model: openai.GPT3Dot5Turbo,
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    openai.ChatMessageRoleUser,
// 					Content: prompt,
// 				},
// 			},
// 		},
// 	)

// 	if err != nil {
// 		return "Error occurred"
// 	}

// 	if len(resp.Choices) > 0 {
// 		return resp.Choices[0].Message.Content
// 	}
// 	return "No reply generated"
// }

// sqlite3 messages.db
