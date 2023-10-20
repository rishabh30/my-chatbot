package messagehandlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"my_chatbot/interfaces"
	"my_chatbot/promptCollection"
	"my_chatbot/repositories"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type MessageHandlers struct {
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // Limit the data size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	msg := r.FormValue("message")
	userid := r.FormValue("userid")
	file, header, err := r.FormFile("file")

	if (err != nil && msg == "") || userid == "" {
		http.Error(w, "Both message and file are missing or userid is missing", http.StatusBadRequest)
		return
	}

	var dstPath string

	if file != nil {
		defer file.Close()

		dstPath = path.Join("storage/files", header.Filename)
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
		log.Println("Successfully uploaded the file: ", header.Filename)
	}

	if msg != "" {
		log.Println("Received message: ", msg)
	}

	// TODO: Add previous 3 message to the msg to get more relavant response
	callOpenAI := promptCollection.Openai{}
	responseText := callOpenAI.CallOpenAIRest(msg)
	id, err := strconv.Atoi(userid)
	if err != nil {
		http.Error(w, "Please pass the user id in correct format", http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	messageRepo := &repositories.MessageRepository{} //TODO: Decouple this piece of code from this method
	var repos interfaces.IMessageRepository
	repos = messageRepo
	repos.StoreMessage(ctxWithTimeout, id, msg, responseText, time.Now(), dstPath)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseText))
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("userid")
	offsetParam := r.URL.Query().Get("offset")
	limitParam := r.URL.Query().Get("limit")

	if userIDParam == "" {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)

	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	messageRepo := &repositories.MessageRepository{} //TODO: Decouple this piece of code from this method
	var repos interfaces.IMessageRepository
	repos = messageRepo
	messages, _ := repos.GetMessage(ctxWithTimeout, userID, limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
