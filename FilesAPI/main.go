package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/streadway/amqp"
)

const (
	URL  = "http://localhost"
	PORT = ":8080"
)

func myHtmlForm(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("index.html")
	url := URL + PORT + "/upload"
	tmpl.Execute(w, url)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File\n")

	//parse input, type multipart/form-data
	r.ParseMultipartForm(10 << 20)

	// retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")

	if err != nil {
		fmt.Println("Error Retrieving file from form-data")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//write temporary file on our server
	tempFile, err := ioutil.TempFile("temp-images", "*.png")

	if err != nil {
		fmt.Println(err)
		return
	}

	info_id := tempFile.Name()

	fmt.Printf("ID name after uploading: %+v\n", info_id)

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	//return whether or not this has been successful
	fmt.Fprintf(w, "Successfully Uploaded File\n")

	producer(info_id)

}

func setupRoutes() {
	http.HandleFunc("/", myHtmlForm)
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(PORT, nil)

}

func producer(path string) {

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Wrong path. No file")
		panic(err)
	}

	arr_dir_id := strings.Split(path, "/")
	id := arr_dir_id[1]

	//The Dial function connects to a server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	q, err := ch.QueueDeclare(
		"Test1",
		false,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	fmt.Println(q)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		fmt.Println(err)
	}

	// attempt to publish a message to the queue!
	err = ch.Publish(
		"",
		"Test1",
		false,
		false,
		amqp.Publishing{
			ContentType: "image/jpeg",
			MessageId:   id,
			Body:        contents,
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")
}

func main() {
	fmt.Println("Go File Upload")
	setupRoutes()

}
