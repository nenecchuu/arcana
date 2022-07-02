# httpclient package

## How to generate mock file using go mock 
1. you can create your own .go file to store mockgen command
2. use this command as example if you want to generate mocks from local package
```
//go:generate mockgen -source=some/relative/source/path.go -destination=some/relative/destination/path.go -package=example HTTPClient
```
3. incase you need to generate the mocks from different repo (say this one only used for dependency in go mod), you can use this template 
```
//go:generate mockgen -destination=some/destination/path.go -package=destination_package github.com/nenecchuu/arcana/pkg/httpclient <interface1, interface2, ...>
```

4. you can check command example in httpclient/client.go