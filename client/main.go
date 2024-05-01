package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	FileServerUrl = os.Getenv("FILE_SERVER_URL")
	os.Mkdir("downloads", os.ModePerm)

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide a command")
		return
	}

	command := args[0]

	service := NewFileUploadService()

	switch command {
	case "upload":

		if len(args) != 2 {
			fmt.Println("Invalid number of arguments")
			return
		}

		dir := args[1]
		key, err := service.UploadFiles(dir)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Store Key: %v\n", key)

	case "get":
		if len(args) != 3 {
			fmt.Println("Invalid number of arguments")
			return
		}

		key := args[1]
		numberStr := args[2]

		number, err := strconv.Atoi(numberStr)
		if err != nil {
			panic(err)
		}

		_, _, err = service.GetFile(key, number)
		if err != nil {
			panic(err)
		}

		fmt.Printf("File %v downloaded and verified\n", number)

	case "demonstration":
		if len(args) != 2 {
			fmt.Println("Invalid number of arguments")
			return
		}

		fmt.Printf("This is a demonstration of the file upload and download service\n")
		fmt.Printf("First we will upload the files from the directory '%v' to the server\n", args[1])
		fmt.Printf("The client will receive Store Key from the server and store Merkle Root of the set of files to local file system\n")
		dir := args[1]
		key, err := service.UploadFiles(dir)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Store Key: %v\n", key)

		fmt.Printf("Now we will download the files from the server via Store Key and file index\n")
		fmt.Printf("Files will be automatically verified using the Proof from the server and local Merkle Root\n")

		getFile(service, key, 0)
		getFile(service, key, 1)
		getFile(service, key, 2)
		getFile(service, key, 3)
		getFile(service, key, 4)
		getFile(service, key, 5)
		getFile(service, key, 6)

	default:
		panic("Invalid command")
	}

	//test()
}

func getFile(service *FileUploadService, key string, number int) {
	_, name, err := service.GetFile(key, number)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File %v (%v) was downloaded into 'downloads' folder and verified\n", number, name)
}

func test() {
	service := NewFileUploadService()
	key, err := service.UploadFiles("files")
	if err != nil {
		panic(err)
	}

	_, _, err = service.GetFile(key, 1)
	if err != nil {
		panic(err)
	}
	_, _, err = service.GetFile(key, 2)
	if err != nil {
		panic(err)
	}

	_, _, err = service.GetFile(key, 3)
	if err != nil {
		panic(err)
	}
}

var FileServerUrl string
