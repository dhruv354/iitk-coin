FROM golang:1.13


#Set the current working directory inside the container

WORKDIR $GOPATH/go/src/github.com/dhruv354/iitk-coin

#commands
copy go.mod .
copy go.sum .
RUN go mod download

#copy everything 
COPY . .

#BUILD Executable 
RUN go build

EXPOSE 8080

#run executable
CMD ["./iitk-coin"]
