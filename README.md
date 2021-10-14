# go-sfomuseum-airfield

Go package for working with airfield-related activities at SFO Museum (airlines, aircraft, airports).

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-airfield.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-airfield)

Documentation is incomplete at this time.

## A note about "lookups"

As of this writing (October, 2021) most of the code in this package is focused around the code for "looking up" a record using one or more identifiers (for example an `ICAO` code or a SFO Museum database primary key) and/or to find the "current" instance of a given identifier (it turns out that IATA codes were re-used in the past). In all cases we're trying to resolve a given identifier back to a specific Who's On First ID for that "thing" in time and the code has been broken out in the topic and source specific subpackages.

All of those subpackages implement a generic `Lookup` interface which looks like this:

```
type Lookup interface {
	Find(context.Context, string) ([]interface{}, error)
	Append(context.Context, interface{}) error
}
```

There is an equivalent interface in the [go-sfomuseum-architecture](https://github.com/sfomuseum/go-sfomuseum-architecture) package which has all the same concerns as this package but is specific to the architectural elements at SFO. In time it would be best if both packages (`go-sfomuseum-architecture` and `go-sfomuseum-airfield`) shared a common "lookup" package/interface but that hasn't happened yet. There were past efforts around this idea in the [go-lookup*](https://github.com/search?q=org%3Asfomuseum+go-lookup) packages which have now been deprecated. If there is going to be a single common "lookup" interface package then the it will happen in a `github.com/sfomuseum/go-lookup/v2` namespace.

Which means that there is _a lot of duplicate code_ to implement functionality around the basic model in order to accomodate the different interfaces because an `icao.Aircraft` is not the same as a `sfomuseum.Aircraft` and a `sfomuseum.Aircraft` won't be the same as a `sfomuseum.Airport`. Things in the `sfomuseum` namespace tend to be more alike than not but there are always going to be edge-cases so the decision has been to suffer duplicate code (in different subpackages) rather than trying to shoehorn different classes of "things" in to a single data structure.

One alternative approach would be to adopt the [GoCloud `As` model](https://gocloud.dev/concepts/as/) but there is a sufficient level of indirection that I haven't completely wrapped my head around so it's still just an idea for now.

It's not great. It's just what we're doing today. The goal right now is to expect a certain amount of "rinse and repeat" in the short term while aiming to make each cycle shorter than the last.

### Data storage for "lookups"

The default data storage layer for lookup is an in-memory `sync.Map`. This works well for most cases and enforces a degree of moderation around the size of lookup tables. Another approach would be to use the [philippgille/gokv](https://github.com/philippgille/gokv) package (or equivalent) which is a simple interface with multiple storage backends. TBD..

## See also

* https://github.com/sfomuseum-data/sfomuseum-data-aircraft
* https://github.com/sfomuseum-data/sfomuseum-data-enterprise
* https://github.com/sfomuseum-data/sfomuseum-data-whosonfirst
