Learning gRPC in Go from this very good course:

https://www.youtube.com/playlist?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf
https://github.com/techschool/pcbook-go

The YouTube course uses the deprecated

  `github.com/golang/protobuf`

but the GitHub repo uses the recommended

  `google.golang.org/protobuf`

In some cases I've slightly changed the package structure from the original course (for example,
all the "storage" code is in its own `storages` package).
