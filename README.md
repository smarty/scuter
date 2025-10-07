# github.com/mdw-go/scuter

> What is a 'scuter'?

A combination of the words "scute" and "scooter".

> Why those words?

Well, this project is a very thin (small) HTTP 'shell' (think ["functional core / imperative shell"](https://www.destroyallsoftware.com/screencasts/catalog/functional-core-imperative-shell), by Gary Bernhardt). A ['scute'](https://en.wikipedia.org/wiki/Scute) is a part or section of, say, a turtle shell. 

Also, this project is a barebones successor to [detour](https://github.com/smarty-archives/detour) and [shuttle](https://github.com/smarty/shuttle), so I wanted a name that fit in with the 'transportation nouns' theme, so something minimal like ["scooter"](https://en.wikipedia.org/wiki/Scooter_(motorcycle)).

`scute + scooter = scuter` 

## Example

The provided example shows how to compose an entire HTTP use case, which:

1. Retrieves data from the HTTP request to form a command,
2. Sends the command to the application layer,
3. Interprets the command's results, and
4. Sends an appropriate HTTP response.

It also demonstrates how you might use the provided `Pool` type to minimize allocations (via `*sync.Pool`). 