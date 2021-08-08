package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/streadway/amqp"
)

const (
	IMAGES_PATH = "../temp_images/"
	BROKER_URL  = "amqp://guest:guest@localhost:5672/"
)

func resizeImage(path string, id string) {

	images_path := fmt.Sprintf("%s%s", path, id)

	file, err := os.Open(images_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	imageData, imageExt, err := image.Decode(file)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("File extension: ", imageExt)

	dstImage128 := imaging.Resize(imageData, 128, 128, imaging.Lanczos)

	out, _ := os.Create(images_path)
	defer out.Close()

	switch imageExt {

	case "jpg", "jpeg":
		err = jpeg.Encode(out, dstImage128, nil)
	case "gif":
		err = gif.Encode(out, dstImage128, nil)
	default:
		err = png.Encode(out, dstImage128)

	}

	if err != nil {
		log.Println(err)
	}

}

func consumer() {
	conn, err := amqp.Dial(BROKER_URL)
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
		"ResizeImage",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Recieved Message: %s\n", d.Body)
			file_id := string(d.Body)
			resizeImage(IMAGES_PATH, file_id)
		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever

}

func main() {

	consumer()

}
