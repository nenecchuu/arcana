# response

`response` is a wrapper that is used for wrapping response for http response.

## Usage Example
for normal case
```
    r := response.NewJSONResponse().
		WithCode("SUCCESS").
		WithData(u).
		WithMessage("some message").
		WithResult(map[string]interface{}{
			"amount": 1,
		}).
		WithLatency(now)
```

if we want to set up error message on response 
```
    r := response.NewJSONResponse().
		WithErrorString("some error")
```

you can find some example in example/response_example