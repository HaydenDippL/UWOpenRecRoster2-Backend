# DOCS

This is a general outline of the code structure.

---

This project is split up into four files: `main.go`, `logging.go`, and `schedule.go`. `main.go` is the entry point of the backend. `logging.go` logs calls to this backend in a TiDB database. `memo.go` memoizes recent calls to the UW-Backend.

## `main.go`

`main.go` contains the single endpoint defined in the README.md. We expect to get the variables `date`, `gyms`, and `facilities`. Is we do not receive one of these, or we receive an invalid request, we send a `400: BAD REQUEST ERROR` back to the client. 

Upon receiving a good request, we will first check to see if we have cached the results in `memo.js`. If we have it memoized in the last hour, use that instead. If we don't have it memoized we will call call the UW-Recwell servers to receive the schedules in the `schedule.go` file.

## `memo.go`

`memo.go` is repsonsible for memoizes/caching results received from `schedule.go`. It has three primary functions. It will be able to check if a schedule has been memoized, it will be able to get a memoized schedule, and it will be able to memoize a schedule (stringified json fromat from README).

Will use a cron job to delete all schedules at the end of a day. We only memoize a month before and after todays date.

## `logging.go`

`logging.go` will log user activity to a remote TiDB database.

## `schedule.go`

`schedule.go` actually fetches and parses the schedules from the UW-Recwell servers. It will format the schedules into the form described in the README.md.