# cloudevents-superfluous-error

```
go get github.com/cloudevents/sdk-go/v2
go get github.com/gorilla/mux
```

```
go build receiver_broken.go
go build receiver_broken_gorilla.go
go build receiver_working.go
go build sender.go
```

Run `./receiver_broken`

In another terminal, run  `./sender` , which send 64 events and queries the receiver on the number of received events

```
2020/06/23 18:49:59 count = 59

```

Anything less than 64 signifies an error.

Compare to `receiver_working`, which always works fine.

