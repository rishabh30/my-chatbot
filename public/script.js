
// Select the chatbot container element
const chatbotContainer = document.querySelector('#chatbot-container');

// Create a function to add a message to the chatbot
function addMessage(message, sender) {
    // Create a new message element
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.classList.add(sender);

    // Create a new message text element
    const messageTextElement = document.createElement('div');
    messageTextElement.classList.add('message-text');
    messageTextElement.textContent = message;

    // Add the message text element to the message element
    messageElement.appendChild(messageTextElement);

    // Add the message element to the chatbot container
    chatbotContainer.appendChild(messageElement);
}

// Create a function to handle user input
function handleUserInput(event) {
    // Prevent the default form submission behavior
    event.preventDefault();

    // Get the user input value
    const userInput = event.target.elements.userInput.value;

    // Add the user input to the chatbot
    addMessage(userInput, 'user');

    // Clear the user input field
    event.target.elements.userInput.value = '';

    let formData = new FormData();

    // For regular fields:
    formData.append('message', userInput);
    formData.append('userid', 1);

    // For file field:
    let fileInput = document.getElementById('myFileInput');
    if (fileInput.files.length > 0) { // Ensure a file was selected
        let file = fileInput.files[0];
        formData.append('field3', file);
    }

    // Send the user input to the chatbot backend
    // (replace this with your own backend API call)
    fetch('/reply', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        // Add the chatbot response to the chatbot
        addMessage(data.message, 'chatbot');
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

// Add an event listener to the chatbot form
const chatbotForm = document.querySelector('#chatbot-form');
chatbotForm.addEventListener('submit', handleUserInput);
