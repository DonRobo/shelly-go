# shelly-go
shelly-go is an unofficial native go client for the [Shelly Gen2 API](https://shelly-api-docs.shelly.cloud/gen2/). It's forked from jcodybaker's excellent [go-shelly](https://github.com/jcodybaker/go-shelly) library. However with the dependency to [mos](github.com/mongoose-os/mos/) completely removed. As many of its transitive dependencies were causing issues for me.

Because of this change, this is not quite a drop in replacement, but it should be simple enough to make the required modifications.

## Maturity
This library is mostly for my own personal amusement and no guarantees or promises about stability, functionality, documentation or test coverage are made.

## Usage
In most cases, using the library is as simple building an Resty client and calling .Do().  See the [go doc](https://pkg.go.dev/github.com/DonRobo/shelly-go), the [Resty doc](https://resty.dev/) and [Shelly Gen2 API documentation](https://shelly-api-docs.shelly.cloud/gen2/).
```
import (
	"context"
	"fmt"
    "log"

	shelly "github.com/DonRobo/shelly-go"
)

const (
    shellyAddress = "http://192.168.1.20"
)

func main() {
	client := resty.New()
	client.SetBaseURL(shellyAddress)
	defer client.Close()

	req := &shelly.ShellyGetStatusRequest{}
	statusResp, _, err := req.Do(client)

    if err != nil {
        log.Fatalf("querying device status: %v", err)
    }
    for i, sw := range statusResp.Switches {
        fmt.Printf("Switch %d status: %v\n", i, *sw.Output)
    }
}
```

## Code generation

The typed config client for many components is generated from the official
[Shelly Gen2+ API docs](https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/)
rather than written by hand, so coverage tracks the documentation.

The generator (`./cmd/gen`, package `./gen`) runs in three stages:

1. **fetch** — download the component doc pages into a local cache under
   `gen/cache/` (gitignored; network is only used for pages not already cached);
2. **parse** — turn the cached HTML into `gen/spec.json` (committed), the
   machine-readable intermediate representation and source of truth for generation;
3. **emit** — generate `<component>_config_gen.go` for every component that does
   not already have a hand-written config type (so generation only fills gaps).

Regenerate with:

```
go generate ./...          # uses the committed cache
go run ./cmd/gen -refresh   # re-download the docs first
```

Per-component corrections the docs can't express generically live in
`gen/exceptions.go`, kept separate from the generic parse/emit code.

Integration tests talk to a real device and are gated behind a build tag:
`go test -tags integration ./...`.

## Legal

### Liability
By downloading/using this open-source library you indemnify the authors of this project from liability to the fullest extent permitted by law. There are real risks of fire, electricution, and/or bodily injury/death with these devices and connected equipment. Errors, bugs, or misinformation may exist in this software which cause your device and attached equipment to function in unexpected and dangerous ways. That risk is your responsibility. 

### License
This fork contains modifications by [DonRobo](https://github.com/DonRobo), based on the original work by J Cody Baker.  
Original copyright [J Cody Baker](cody@codybaker.com)  
Licensed under the [MIT license](LICENSE.md).