# iitk-coin

## It is a Pseudo currency system Project at IIT Kanpur where students can transact among themselves with this currency

### Programming language Used - Golang and database - sqlite3

### Directory  conventions and other instructions
* All the code should lie within `path=$GOPATH/go/src/github.com/your username`
* Initialize your Go module by running `go mod init` which creates go.mod file which contains all the dependecies

### Currently using Two tables in tha database one is `User` for storing signup details and one is `Userdata` for storing coins of that user
### Routes
* signup and login route for authentication purposes
* getcoins, addcoins and transfercoins routes for getting user current coins, awarding coins to the user say for participation and transfering coins between two different users