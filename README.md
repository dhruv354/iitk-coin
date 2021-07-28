# IITK Coin
## SnT Project 2021, Programming Club 

This repository contains the code for the IITK Coin project done so far.

### Relevant Links

- [Proposal](https://docs.google.com/document/d/1Jm2pImnVrgxi7Qu-DvHG4wdYIYLJ4IYw4N9uCU-z_t0/edit?usp=sharing)
- [Poster](https://www.canva.com/design/DAElW0YxGAg/XqLWvQp10V5lCzMTEWRBCQ/view?utm_content=DAElW0YxGAg&utm_campaign=designshare&utm_medium=link&utm_source=sharebutton)
- [Midterm Evaluation presentation](https://docs.google.com/presentation/d/1kriN-7A3v1RlXUDL5NETX3roJKRMJInptkWofIxY8dg/edit?usp=sharing)
- [Midterm Documentation](https://docs.google.com/document/d/1bvOWH4k0U-l2pQ1jLWIDzOkJ2wbHNW4jJw7tMWkUV6o/edit?usp=sharing)

## Table Of Content
- [Development Environment](#development-environment)
- [Directory Structure](#directory-structure)
- [Usage](#usage)
- [Endpoints](#endpoints)

## Development Environment

```bash
- go version: go1.16.4 linux/amd64   
- OS: ubuntu-20.04 LTS   
- text editor: VSCode    	
- terminal: ubuntu terminal 
```

## Directory Structure
```
.
├── README.md
├── Handlers
│   └── handler.go
├── sqlite3_func
│   └── sqlite3Func.go
├── utilities
│   └── utilities.go
├── go.mod
├── go.sum
├── Student_info.db
├── main.go


## Usage
```bash
cd $GOPATH/src/github.com/<username>
git clone https://github.com/dhruv354/iitk-coin.git
cd repo
go run main.go     
#, or build the program and run the executable
go build
./iitk-coin
```

Output should look like

```
created my database
User table created or not altered if already created
user coins table created
Serving at 8080
```

## Endpoints
POST requests take place via `JSON` requests. A typical usage would look like

```bash

```

- `/login` : `POST`
```json
{"rollno":"<rollno>", "password":"<password>"}
```

- `/signup` : `POST`
```json
{"name":"<username>","rollno":"<user rollno>", "password":"<password>","batch":"<user batch>"}
```

- `/logout` : `POST`
```json

```
- `/redeemcoins` : `POST`
```json
  {"coin":"<How much coin want to redeem>", "item":"<item name>"}
```

- `/itemredeem` : `POST`
```json
{"item":"<Id of redeem request>", "coins":"<price of that item>"}
```

- `/adminApproval` : `POST`
```json
{"item":"<Id of redeem request>"}
```

- `/addcoins` : `POST`
```json
{"coins":"<Coins to reward>", "rollno":"<rollno to which coin will be rewarded>"}
```

- `/transfercoins` : `POST`
```json
{"coins":"<Coins to transfer>", "rollno":"<Whom to transfer>"}
```

GET requests:

- `/secretpage` : `GET`
```bash
```

- `/getcoins` : `GET`
```bash
curl http://localhost:8080/getcoin
```

## Models

-  UserData
```go
	Name     string `json:"name"`
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
	Batch    int    `json:"batch"`
```

- Claims
```go
	Username string `json:"username"`
	Rollno   int    `json:"rollno"`
	jwt.StandardClaims
```

- UserCoins
```go
  Rollno int `json:"rollno"`
	Coins  int `json:"coins"`
```

- TransferBWUsers
```go
  Rollno int `json:"rollno"`
	Coins  int `json:"coins"`
```

- ItemRedeem
```go
	Item  string `json:"item"` 
	Coins int    `json:"coins"` 
```


- RequestId
```go
		Id int `json:"id"`
```
