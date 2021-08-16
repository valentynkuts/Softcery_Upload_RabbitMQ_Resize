# Upload_RabbitMQ_Resize

```bash
docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```

It consists of two microservices, one responsible for file uploading and retrieval (Files API), and another responsible for file processing (Processing API).

The Files API microservice expose an HTTP endpoint for uploading files. Once a file is uploaded, it is saved to the file system, and its ID is sent to the Processing API via a RabbitMQ queue. 

The Processing API microservice provides functionality for processing and optimizing images. It accepts file ID from a RabbitMQ queue, reduces the image size and overwrites the file in the file system.
