数据类型

* Counter
* Gauge
* Histogram
* Summary

Status Code

* 400
* 422 Unprocessable Entity (expression can't be executed)
* 503 Service Unavailable


Query

single instant

query
time
timeout


range of time

# Response

{
  "status": "success" | "error",
  "data": <data>,

  // Only set if status is "error". The data field may still hold
  // additional data.
  "errorType": "<string>",
  "error": "<string>",

  // Only if there were warnings while executing the request.
  // There will still be data in the data field.
  "warnings": ["<string>"]
}
