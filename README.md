# check-deps

[![Go Report Card](https://goreportcard.com/badge/github.com/nivl/check-deps)](https://goreportcard.com/report/github.com/nivl/check-deps)

Check dep is a program that parses the `go.mod` a of package (current
directory) and prints all the modules that need to be updated, or that haven't
received any updates in the last 6 months (they might no longer be maintained).

## Install

`go get -u github.com/Nivl/check-deps`

## usage

`check-deps [flags]`

| Flag         | Description                                             |
| ------------ | ------------------------------------------------------- |
| --old        | check for modules without updates for the last 6 months |
| --ignore, -i | coma separated list of packages to ignore               |
| --indirect   | check indirect modules                                  |

## Output

```
❯ check-deps --old --indirect -i github.com/Nivl
+-------------------------------------+------------------------------------+------------------------------------+----------+
|               MODULE                |          CURRENT VERSION           |            NEW VERSION             | INDIRECT |
+-------------------------------------+------------------------------------+------------------------------------+----------+
| golang.org/x/crypto                 | v0.0.0-20181001203147-e3636079e1a4 | v0.0.0-20190426145343-a29dc8fdc734 |  false   |
| golang.org/x/lint                   | v0.0.0-20180702182130-06c8688daad7 | v0.0.0-20190409202823-959b441ac422 |   true   |
| golang.org/x/net                    | v0.0.0-20181005035420-146acd28ed58 | v0.0.0-20190501004415-9ce7a6920f09 |  false   |
| golang.org/x/oauth2                 | v0.0.0-20181003184128-c57b0facaced | v0.0.0-20190402181905-9f3314589c9a |  false   |
| golang.org/x/sync                   | v0.0.0-20180314180146-1d60e4601c6f | v0.0.0-20190423024810-112230192c58 |   true   |
| golang.org/x/sys                    | v0.0.0-20181005133103-4497e2df6f9e | v0.0.0-20190429190828-d89cdac9e872 |  false   |
| golang.org/x/text                   |               v0.3.0               |               v0.3.2               |  false   |
| golang.org/x/tools                  | v0.0.0-20180828015842-6cd1fcedba52 | v0.0.0-20190501045030-23463209683d |   true   |
| google.golang.org/api               | v0.0.0-20181007000908-c21459d81882 |               v0.4.0               |  false   |
| google.golang.org/appengine         |               v1.2.0               |               v1.5.0               |  false   |
| google.golang.org/genproto          | v0.0.0-20181004005441-af9cb2a35e7f | v0.0.0-20190425155659-357c62f0e4bb |  false   |
| google.golang.org/grpc              |              v1.15.0               |              v1.20.1               |  false   |
| honnef.co/go/tools                  | v0.0.0-20180728063816-88497007e858 | v0.0.0-20190418001031-e561f6794a2a |   true   |
+-------------------------------------+------------------------------------+------------------------------------+----------+

+--------------------------------------------------+----------------------------+----------+
|                      MODULE                      |        LAST UPDATE         | INDIRECT |
+--------------------------------------------------+----------------------------+----------+
| github.com/bsphere/le_go                         | 26 months ago (2017/02/15) |  false   |
| github.com/client9/misspell                      | 13 months ago (2018/03/09) |   true   |
| github.com/davecgh/go-spew                       | 14 months ago (2018/02/21) |  false   |
| github.com/dchest/uniuri                         | 39 months ago (2016/02/12) |  false   |
| github.com/golang/glog                           | 39 months ago (2016/01/26) |   true   |
| github.com/gorilla/context                       | 32 months ago (2016/08/17) |  false   |
| github.com/gorilla/handlers                      | 9 months ago (2018/07/27)  |  false   |
| github.com/kelseyhightower/envconfig             | 27 months ago (2017/01/24) |  false   |
| github.com/kisielk/gotool                        | 14 months ago (2018/02/21) |   true   |
| github.com/matttproud/golang_protobuf_extensions | 36 months ago (2016/04/24) |   true   |
| github.com/pmezard/go-difflib                    | 40 months ago (2016/01/10) |  false   |
| github.com/rainycape/unidecode                   | 44 months ago (2015/09/07) |  false   |
| github.com/satori/go.uuid                        | 16 months ago (2018/01/03) |  false   |
| github.com/sendgrid/rest                         | 12 months ago (2018/04/09) |  false   |
| github.com/sendgrid/sendgrid-go                  | 22 months ago (2017/07/04) |  false   |
+--------------------------------------------------+----------------------------+----------+
```
