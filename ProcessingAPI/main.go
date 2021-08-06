package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/streadway/amqp"
)

func processing(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Processing\n")

	consumer()
}

func setupRoutes() {
	http.HandleFunc("/", processing)
	http.ListenAndServe(":8081", nil)

}

func serveFrames(imgByte []byte, id string) {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}

	dstImage128 := imaging.Resize(img, 128, 128, imaging.Lanczos)

	path := fmt.Sprintf("photo/%s", id)
	fmt.Println(path)
	out, _ := os.Create(path)
	defer out.Close()

	err = png.Encode(out, dstImage128)

	if err != nil {
		log.Println(err)
	}

}

func consumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	msgs, err := ch.Consume(
		"Test1",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	fmt.Println(msgs)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//fmt.Printf("Recieved Message: %s\n", d.Body)
			//fmt.Println(d.Body)
			fmt.Println(d.MessageId)
			id := string(d.MessageId)
			serveFrames(d.Body, id)
		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever

}
func main() {
	fmt.Println("Processing")
	setupRoutes()

}
