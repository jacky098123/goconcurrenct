# Requirement
2 important requirements which impact the design a lot

## Half sync and half async. 
This is important when calling from external. For example: In STG, a receivable creation request can get all context to help engineer what is going on, but In production, if there are too many data, the endpoint can't wait infinitely, but return a requestID to let engineer check the Kibana log.

### Parent goroutine need to know all the details of child goroutine
* The child goroutine returned an error, parent goroutine can return an resonable error
* The child is running in Async mode, can tell caller, this request become an async call, can check Kibana log for detail
* The parent goroutine is canceled by synchronize control
* There is a panic

the API design should be closure return an error, or even more data
```
backgroundFunc = func() error {
  // do something
  if err {
    return err
  }
  return nil
}

backgroundFuncWithResult = func() (interface{}, error) {
}

goWithRecovery(ctx, backgroundFunc, timeout)
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
