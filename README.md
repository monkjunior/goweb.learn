# Teaching myself web development with Go

```bash
#Run the project
go get github.com/pilu/fresh
fresh
```

- How fast can i init this project ? 8m:20

- Web app usually does not need to be go getable :D -> Just name repo as the domain you want to host.

- When you do a GET request you dont have your body of request.

- Serve mux acts like a router of HTTP reuqests.

- Dynamic reloading for local development: 

```bash
go get github.com/pilu/fresh
cd path/go/myapp
fresh
```
- Should explicit set the Content-Type header if you can.

- HTTP status codes
[net/http header write](https://pkg.go.dev/net/http#Header.Write)

> If WriteHeader has not yet been called, Write calls
> WriteHeader(http.StatusOK) before writing the data.

- ServeMux: pattern matching and why we should not use this mux?

    - Exactly match, care of your traling slask */*

    - Longest pattern match

    - "/"  matches everything

    - No dependencies

    - TOO simple

    - Useful for simple testing or set up a middleware for different router, domain-based handling.

    - Missing features: regex matching pattern, named url params, dynamic url.

- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)

    - Simple, fast, Mem & CPU friendly

    - If your application needs to serve HTTP files, this is the best choice

    - If your application spends most of the time dealing with DB, read things from disk and so on, this router should not be your primary concern

    - Named parms -> Dynamic URL

    - Support HTTP method for route matching

 - [gorilla web toolkit](https://www.gorillatoolkit.org/)

    - Not suitable for newbie

    - Support {named} URL params, minor pattern

    - Even gorilla support domain matching but the author suggest using default mux for this job.

    - Also support HTTP methods

    - Can easily switch between gorilla mux and default build-in mux

    - Auto handle 404


## Template

Template supports putting dynamic data into a text file.

Logic should be handled before you go to the template handling.

Go provides text/template and html/template package. They use a similar interface.

html/template package:

    - Auto do HTML encoding to prevent code injection


## MVC: Model-View-Controller

How the web request come and travel along these components ?

View: rendering data. The data could be in any kind: json, html, etc. View should least logic as possible.

    - Create reusable layouts

Controller: middle man, router

    - Control the flow of our codes

    - Should not have a ton of logic in it

    - The way controllers use to talk to other components is like "Hey, i have some stuff to do, pls do this for me!". They does not do anything much.

Model: connect to database, connect to other APIs

    - Interact with our data. Data can be DBs, APIs, Files, S3, etc.

    - Can do validate data here.

Not everything need to design with MVC.

## Bootstrap

A framework for HTML, CSS and JS

## User controller

An action = Handler function.

    - Take a req, kick off some logic and then respone

Restful design.

Use method instead of function let you access to object data without public a bunch of global variables.

