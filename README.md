# GOIODI

GOIODI is an open-source responsive collaborative dictionary web application made with the Ionic framework, AngularJS, the Golang language and based on a MongoDB NoSQL database. The main advantages of GOIOD:

  - Open source
  - Has a responsive design and can be used on mobile devices (an application will be created for Android)
  - Collaborative (social dictionary for small and big communities)
  - Can be translated to several languages (only French and English are supported for the moment)

### Current version
1.2

### Used technologies

GOIODI is powererd with a number of open source projects:

* [Ionic] - An HTML5 Hybrid Mobile App Framework.
* [AngularJS] - An HTML/JavaScript MVW Framework.
* [Golang] - A modern programming language.
* [MongoDB] - A NoSQL database.
* [Angular-Translate] - An AngularJS translation library.

And of course GOIODI itself is open source with an MIT license, you can use it for commercial/non-commercial products for free.

### Installation

Download latest package lists:
```sh
sudo apt-get update
```
Install Golang 1.6 (https://golang.org/doc/install?download=go1.6.linux-amd64.tar.gz)

Clone the project repository or download the ZIP file from this page then build the Golang server:
```sh
go build
```
Install and run MongoDB 3.2 (https://docs.mongodb.org/manual/installation/)
```sh
sudo service mongod start
```
If you want to populate the database, use the provided script in the project's root directory.
```sh
./populateDB_TLS
```

### Running the HTTP server

You need to generate TLS certificates in the project repository's root directory (the certificate files must be placed under /opt):
```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host="localhost"
```
or
```sh
go run /usr/lib/go/src/pkg/crypto/tls/generate_cert.go --host="localhost"
```
Then, run the back-end server in secure mode:
```sh
sudo ./goiodi
```

### Todos

 - Add word highlighting when user searches a word
 - Add word vocal search
 - Add word audio reading
 - Add user related features:
    - Add word ranking (with a 5 star system for example)
    - Add word consultation statistics
    - Add "Word of the day" feature to make the user learn new words

License
----
MIT
