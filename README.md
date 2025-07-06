# Whats-Valid

Whats-Valid is a simple web application that allows you to check if a phone number is registered on WhatsApp. It features a Go backend that connects to WhatsApp and a clean, responsive frontend built with Svelte.

## 🚀 Features

- **WhatsApp Number Validation**: Checks if a phone number has an active WhatsApp account.
- **QR Code Login**: Easily connect to WhatsApp by scanning a QR code with your phone.
- **Web Interface**: A user-friendly web UI to enter and check phone numbers.
- **All-in-One Binary**: The Go backend serves the frontend, making it a single, self-contained application.

## 🛠️ Technologies Used

- **Backend**:
  - [Go](https://golang.org/)
  - [Whatsmeow](https://github.com/maunium/go-whatsmeow): A Go library for the WhatsApp Web API.
  - [go-sqlite3](https://github.com/mattn/go-sqlite3): For storing session data.
- **Frontend**:
  - [SvelteKit](https://kit.svelte.dev/)
  - [Tailwind CSS](https://tailwindcss.com/)
  - [Vite](https://vitejs.dev/)

## 🏁 Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Make sure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [Node.js](https://nodejs.org/) (version 18 or higher)
- [pnpm](https://pnpm.io/installation) (or you can use `npm` or `yarn`)

### Installation & Running

1.  **Clone the repository:**

    ```bash
    git clone [https://github.com/wiscaksono/whats-valid.git](https://github.com/wiscaksono/whats-valid.git)
    cd whats-valid
    ```

2.  **Set up the Frontend:**
    Navigate to the `frontend` directory and install the necessary dependencies.

    ```bash
    cd frontend
    pnpm install
    ```

    After installation, build the frontend. This will create a `build` directory that the Go backend will serve.

    ```bash
    pnpm run build
    ```

3.  **Run the Backend:**
    Go back to the root directory. Build and run the Go application.
    ```bash
    cd ..
    go build -o whats-valid main.go
    ./whats-valid
    ```
4.  **Login with WhatsApp:**
    The first time you run the application, a **QR code** will be displayed in your terminal. Scan this QR code with the WhatsApp application on your phone (in `Settings > Linked Devices > Link a Device`). Once you've successfully logged in, the server will start. Your session will be saved in `store.db`, so you won't need to scan the QR code every time you start the server.

5.  **Access the Application:**
    Open your web browser and navigate to:
    ```
    http://localhost:3000
    ```
    You should now see the Whats-Valid web interface, ready to check numbers!

## ☁️ API Endpoint

The application exposes a single API endpoint:

### `GET /check`

This endpoint checks if a given number is on WhatsApp.

- **Query Parameter:**
  - `number` (string, required): The phone number to check, including the country code (e.g., `15551234567`).

- **Success Response (200 OK):**

  ```json
  {
    "number": "15551234567",
    "isOnWhatsApp": true,
    "status": "success"
  }
  ```

- **Error Responses:**
  - `400 Bad Request`: If the `number` parameter is missing.
  - `503 Service Unavailable`: If the WhatsApp client is not connected.
  - `500 Internal Server Error`: For other failures during the check.

## 📂 Project Structure

.
├── main.go # The main Go application file
├── go.mod # Go module dependencies
├── .air.toml # Configuration for live-reloading (development)
├── frontend/ # SvelteKit frontend source code
│ ├── src/
│ ├── static/
│ ├── package.json
│ └── svelte.config.js
└── ...
