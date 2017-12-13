# Agenda

![TravisCI](https://travis-ci.org/Ace-0/AgendaService.svg?branch=master)

Agenda is a command line meeting manager.

## Travis Testing

[.travis.yml](https://github.com/ZexinYan/AgendaService/blob/master/.travis.yml)

- Build `service` and `cli`
- Run tests of `service` and `cli`

## Docker Integration

[AgendaService on dockerhub](https://hub.docker.com/r/vinalx/agendaservice/)

```
docker pull vinalx/AgendaService
```

### Start client in docker

```
$ docker run --rm agenda-service client -h
Agenda is a meeting manager based on CLI using cobra library.
It supports different operation on meetings including register, create meeting, query and so on.
It's a cooperation homework assignment for service computing.

...
```

### Start server in docker

```
$ docker run --rm --volume "dbdir:/go/db/" agenda-service server -h
flag needs an argument: 'h' in -h
Usage of /go/bin/service:
  -d, --db string
    	Path to the sqlite db file (default "agenda.db")
  -h, --host string
    	Host of the server listening on (default "0.0.0.0")
  -p, --port value
    	Port of the server listening on (default 8080)
```

## Usage Examples

```
go get -v github.com/ZexinYan/AgendaService/service
go get -v github.com/ZexinYan/AgendaService/service
```

or use docker

### Help

```
$ ./cli
Agenda is a meeting manager based on CLI using cobra library.
It supports different operation on meetings including register, create meeting, query and so on.
It's a cooperation homework assignment for service computing.

Usage:
  Agenda [command]

Available Commands:
  delete      Delete your account.
  help        Help about any command
  info        Show the information of your account.
  list        list Users
  login       Login
  logout      Logout
  register    Register user.

Flags:
  -b, --baseUrl string   Specify the server url. (default "http://0.0.0.0:8080")
  -d, --debug            display log message
  -h, --help             help for Agenda

Use "Agenda [command] --help" for more information about a command.

```


### Create User

```
$ ./cli register -h
You need to provide username and password to register, and the username can't be the same as others.

Usage:
  Agenda register [flags]

Flags:
  -h, --help              help for register
  -m, --mail string       email.
  -p, --password string   Help message for username
  -t, --phone string      Phone
  -u, --user string       Username

Global Flags:
  -d, --debug   display log message
```


```
username: yanzexin
password: 123
mail: yzx9610@outlook.com
phone: 15626411416
```

If the username hasn't been registered, then you will succeed.

```
$ ./cli register -u yanzexin -p 123 -m yzx9610@outlook.com -t 15626411416
Register Succeed
```

However, if the username has been registered, then you will fail.

```
$ ./cli register -u yanzexin -p 123 -m yzx9610@outlook.com -t 15626411416
Conflict
```

### Login

```
$ ./cli login -u yanzexin -p 123
Login Succeed
```


```
$ ./cli login -u yanzexin -p 1234
Forbidden
```

### Log out

```
$ ./cli logout
Log Out Succeed!
```

### Show user's information
If you want to retrieve the information, you have to log in first.

```
$ ./cli info -u yanzexin
Please Login first!
```

You can get the information about yourself...

```
$ ./cli info -u yanzexin
{
	"email": "yzx9610@outlook.com",
	"phone": "15626411416",
	"username": "yanzexin"
}
```

Or?

```
$ ./cli info -u yan
{
	"email": "mail",
	"phone": "1234",
	"username": "yan"
}
```

### Show all users

```
$ ./cli list
{
	"email": "mail",
	"phone": "1234",
	"username": "yan"
}
{
	"email": "yzx9610@outlook.com",
	"phone": "15626411416",
	"username": "yanzexin"
}
```

### Delete user
Actually, we don't recommend you to do that, but if you want...


```
$ ./cli delete -u yanzexin
Delete User Succeed!
```

Then, you will lose all of your data in Agenda.

## Test

### CLI

```
$ go test
Testing Register...
Register Succeed
Testing Login...
Login Succeed
Testing Log out...
Login Succeed
Log Out Succeed!
Testing show information...
Login Succeed
{
        "email": "mail",
        "phone": "1234",
        "username": "yanzexin"
}
Testing show users...
Login Succeed
{
        "email": "mail",
        "phone": "1234",
        "username": "yan"
}
{
        "email": "mail",
        "phone": "1234",
        "username": "yanzexin"
}
Testing delete user...
Login Succeed
Delete User Succeed!
PASS
ok      github.com/AgendaService/cli/vendor/cmd 0.028s

```

### Server

```
$ go test ./service/vendor/model
ok  	homework/sc/AgendaService/service/vendor/model	0.051s
$ go test ./service/vendor/database
ok  	homework/sc/AgendaService/service/vendor/database	0.095s
```