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

## DRY up your code (dont repeat yourself)

- This is error prone.

- Shorter code != easier to understand or maintain

## Persisting your data

You can use anything to persist data:

- CSVs

- DBs, each has its pros and cons that  you need to understand when you consider to pick one.

PostgreSQL:

- Massive scale: billion of users

- Educational resources

    - [Using PostgreSQL with Golang](https://www.calhoun.io/using-postgresql-with-go/)

    - [Codecademy's Learn SQL course](https://www.codecademy.com/learn/learn-sql)

    - [w3sschools's course](https://www.w3schools.com/sql/)

    - [quora](https://www.quora.com/How-do-I-learn-SQL)

[Why we import packages that we dont actually use ?](https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/)

```bash
docker run -d \
    --name goweb-postgres \
    -p 5432:5432 \
    -e POSTGRES_USER=ted \
    -e POSTGRES_PASSWORD=your-password \
    -e POSTGRES_DB=goweb_dev \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v postgresql:/var/lib/postgresql/data \
    postgres
```

Relational datas:

- GORM, ORM stands for Object Relational Mapping

    - Define models

    - AutoMigrate

    - Validate fields both at DB layer and backend level, or even fronend level

    - Break .Where() to smaller chains makes your code more readable.

    - Show relational data: Preload loads the association data in a separate query.

        - This will not use JOIN

## Creating User model

The Controller layer dont need to know what database engine that the model is using.

    - We will create an interface to seperate 2 layers -> A service layer

User's attributes:

    - gorm.Model: id, created_at, updated_at, deleted_at

    - Name

    - Email

We will not put in password field until we make sure that we stored other things correctly.

Before you really understand your code that you are writing, adding more automate tests is not necessary.

    - Testing Patterns ??

But in go, to run test, you need to run
```bash
go test .
```

## AuthN System

This is the most important and sensitive part of your app.

But implement this is not that hard. There are many small but relatively simple steps.

DO NOT deviate from the norms!

Why dont we use a third party package ?

- Every web dev should have a basic understanding of secure auth practises.

- Custome requirements mean you still need to customize most off-the-shelf solutions.

- It could save your time in the long run.

- It could alse save your money over using services like Auth0 or Stormpath.

Customer requirements and basic understanding

- Without understand basic security, it is easy to compromise your system when you make these customizations, so you can not avoid learning this.

Saving time

- Custom our auth is easy.

Do not reinvent the wheel

Always use SSL/TLS in prod

- Password should never go over the wire unencrypted

- Cookie theft

Hashing password

- If you can decrypt the password, you are doing it wrong! You are not hashing but encrypting.

- These two are VERY different.

- Your app should NEVER be able to recreate a password. So how do we verify a user's password when they log in ?

- Instead of storing a password, we store a hash.

    - bcrypt

- Why cann't they just reverse those hashes ?

- Salt and Pepper are techniques that make hackers harder to figure out our passwords.

    - Salting: giving every user password some random strings added to it before we hash it.

        - Diff users can have diff salts.

        - Can be stored in databases.
    
    - Peppers: same idea with salting technique. But if salt is applied to each users, pepper is applied to each application.

        - Not stored in DB but our application (our codes).

        - Not entirely necessary.

## Remembering users

Web servers are stateless means:

- The server handles each request independently.

- Servers will not remember what you did 15 mins ago.

- Pros:

    - Users can talk to different servers each request -> Easy to scale.

    - Server outages dont lose work.

    - Easier to code - each request has all the data if needs.

So how do we remember who a user is then ?

- We dont.

- We let the user tell us who they are every web request!

    - How do users tell us who they are ?

    - What if a user lies to us ?

    - This does not make any sense!!!

What are cookies?

- Data file stored on a user's computer.

- Linked to a specific domain.

- Both user and server can edit the cookies.

What are cookies used for ?

- Session

- JWTs can be stored in cookies

Create our first cookie

- Store an email address