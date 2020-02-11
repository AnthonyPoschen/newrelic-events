# newrelicEvents

GO SDK for posting new relic events

## why?

after needing to use newrelic at scale on production systems it became
noticable the performance impact of the very bloated Newrelic SDK was
having on production systems, so this package was created to keep it simple
and strip back to just event posting mechanism in an attempt to improve performance.

## Usage

to get this package
```sh
go get -u gitlab.com/zanven/newrelicEvents
```
example usage of the package
```golang
func main() {
    nr := newrelicEvents.New("APPLICATION_ID","LICENCE")
    
    myEvent := map[string]interface{}{"var":"value"}

    // RecordEvent is safe to use concurrently
    go nr.RecordEvent("custom_event_type",myEvent)

    // ensure all data is posted before shutdown
    err := nr.Sync()
}
```
