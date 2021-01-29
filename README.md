# newrelicEvents

GO SDK for posting new relic events

https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/introduction-event-api

## why?

after needing to use newrelic at scale on production systems it became
noticable the performance impact of the very bloated Newrelic SDK was
having on production systems, so this package was created to keep it simple
and strip back to just event posting mechanism in an attempt to improve performance.

## Usage

to get this package
```sh
go get -u github.com/zanven42/newrelic-events
```
example usage of the package
```golang
func main() {
    nre := events.New("APPLICATION_ID","LICENCE")
    
    myEvent := map[string]interface{}{"ENV":"DEV","APP":"newrelic-events-hello-world":"var3":"val3"}

    // RecordEvent is safe to use concurrently
    go nre.Record("custom_event_type",myEvent)

    // ensure all data is posted before shutdown
    err := nre.Sync()
}
```
