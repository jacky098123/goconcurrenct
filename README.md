# Requirement
2 important requirements which impact the design a lot

## Half sync and half async. 
This is important for request processing. 
For example: In STG, the workload is small, a receivable creation request can be done with seconds, to return the result in response is very friendly to the engineer.
But in production, the data is big, it will takes over 1 hours sometimes, the endpoint can't wait infinitely, alternative return a requestID in response and let the engineer to check the Kibana log.

### Parent goroutine need to know all the details of child goroutine
* The child goroutine returned an error, parent goroutine can return an reasonable error
* There is a panic in child goroutine, the parent goroutine can catch it and return an error
* The child can run in Async mode after the wait for some time, then this request become an async call, return the requestID
* The parent goroutine is canceled by synchronize control

the API design should be closure return an error, or even more data
```
backgroundFunc = func() error {
  // do something
  if err {
    return err
  }
  return nil
}

retErr := goWithRecovery(ctx, backgroundFunc, timeout)

// or return a value
backgroundFuncWithResult = func() (interface{}, error) {
}
```

## Parent goroutine can control child goroutine
For example: you are running data deletion from RDS periodly, If previous goroutine is still running, then new goroutine will compete with old one.
This control generall is controld by context.

```
backgroundFunc = func() error {
  nexCtx = context.WithTimeout(context.Background(), duration)
  // do something with newCtx
  if err {
    return err
  }
  return nil
}
```

## Other consideration
### Dependencies
When you implement it in your company, it is criticial to make it as simple as possible to simple integration. 2 main considerations is DdataDog, Logging

So the best way is to add a more layer to integrate it with different dependencies
