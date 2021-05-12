# brabbit
A CLI tool for sending rabbit messages to a selected queue / exchange using predefined headers and message bodys saved as JSON.


## Useage

Run `go get github.com/willdot/brabbit` to install it.

You can send messages to either a queue or an exchange (not both).

Create 2 JSON files, one to contain your headers and another which will contain the message body.

For example:

`headers.json`
``` json
{
  "some header" : "value"
}
```

`body.json`
```json
{
  "some field" : "some value"
}
```
Then run the command:

`brabbit  -body=body.json -headers=headers.json` 

For a queue use the `-queue={queue name}` flag and for exchange use the `-exchange={exchange name}` flag

Note: Only header type exchanges are supported at the moment.


You can also send the same message multiple times by using the `-repeat=10` flag
