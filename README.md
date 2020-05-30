Image Store Service in Golang


Run the service with:
go run main.go    

To create new album:
curl --location --request POST 'http://localhost:8000/create/{name}'

To upload image:
Browse 'http://localhost:8000/upload/'

To delete:
curl --location --request POST 'http://localhost:8000/delete/{name}'

To list specific album:
curl --location --request POST 'http://localhost:8000/get/{name}'

To list all:
curl --location --request POST 'http://localhost:8000/getall/'