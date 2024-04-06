# Requirement
2 important requirements which impact the design a lot

## Half sync and half async. 
This is important when calling from external. For example: In STG, a receivable creation request can get all context to help engineer what is going on, but In production, if there are too many data, the endpoint can't wait infinitely, but return a requestID to let engineer check the Kibana log.

## Parent goroutine can control child goroutine
For example: you are running data deletion from RDS periodly, If previous goroutine is still running, then new goroutine will compete with old one.

## Other consideration
### Dependencies
When you implement it in your company, it is criticial to make it as simple as possible to simple integration. 2 main considerations is DdataDog, Logging

So the best way is to add a more layer to integrate it with different dependencies
