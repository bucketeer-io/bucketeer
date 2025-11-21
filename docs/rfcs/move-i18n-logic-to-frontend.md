# Summary

Now, we have part of the i18n implementation on both sides, front and backend. 

Because the i18n implementation in all APIs made the backend code too hard to read, we decided to do this only on the front end.

[Issue](https://github.com/bucketeer-io/bucketeer/issues/1253)

## Scope

There are two parts that have i18n logic in backend:

1. error message
2. domain event

**We focus on only [1. error message] part.**

We will skip [2. domain event] for now for the following reasons:
  1. Lower Priority: Let's develop Analysis Dashboard first to make Bucketeer more useful.
  2. Less Impact: The localization in the Domain Events affects only on Audit Logs and Slack notifications. There's no way to migrate i18 logic to frontend in the Slack notifications.


## New Error Design

Currently, we use GRPC's `ErrorDetail` to return a localized message based on [GRPC's `LocalizedMessage`](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto#L290).

Example: Validation error when GetAccount

https://github.com/bucketeer-io/bucketeer/blob/main/pkg/account/api/validation.go#L624

We're going to design to return by utilizing `ErrorInfo` defined in GRPC's [GRPC's ErrorInfo](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto#L51).


### Response Format

Response Example:
```json
{
  "code": 3,
  "message": "rpc error: code = InvalidArgument desc = account:invalid email",
  "details": [
    {
      "reason": "INVALID",
      "domain": "account.bucketeer.io",
      "metadata": {
        "messageKey": "InvalidArgumentError",
        "email": "email.com",
        "field_1": "APIKey"
      }
    }
  ]
}
```
| Key                     | Explanation                                                                                                | Example                                              |
| :---------------------- | :--------------------------------------------------------------------------------------------------------- | :--------------------------------------------------- |
| reason                  | The reason of the error.                                                                                   | "reason": "INVALID"                                  |
| domain                  | The error domain is typically the registered service name of the tool or product that generates the error. | "domain": "account.bucketeer.io"  // account package |
| metadata                | Additional structured details about this error.                                                            | -                                                    |
| metadata.messageKey     | Key to identify message content.<br>                                                                       | e.g. NotFoundError, InvalidArgumentError             |
| metadata.<key-value(s)> | Additional information to be embedded in the message. Optional.                                            | "email": "email.com", "field_1": "APIKey"            |

For now, `details` has only one element.

### Message Formats

We will move the error message formats in `pkg/locale/localizedata/` to frontend, like under `ui/dashboard/src/@locales`.
They contain both error messages and nouns.

e.g.

```json
// en/backend-errors.json
{
    "NotFoundError": "The requested {{ .field_1 }} cannot be found",
    "InvalidArgumentError": "The argument {{ .field_1 }} is invalid",
    "ExceededMaxError": "The maximum value {{ .field_2 }} for {{ .field_1 }} has been exceeded",
    ...
}

// en/nouns.json (The file is not decided yet)
{
    "APIKey": "API key",
    "OffVariation": "Off variation",
    ...
}
```


## Correction points
### Point1: Add common functions that generate errors
Create `NewError()` for each package.

Ex: account Package
```go
func NewError(status *gstatus.Status, anoterDetailData ...map[string]string) error {
	domain := "account.bucketeer.io"
	var details []*errdetails.ErrorInfo
	var reason string
	var messageKey string
	var metadatas []map[string]string
	if status == statusEmailIsEmpty {
		reason = "INVALID"
		messageKey = "RequiredFieldError"
		metadatas = []map[string]string{
			{
				"messageKey": messageKey,
				"field":      "email",
			},
		}
	} else if status == statusInvalidEmail {
		reason = "INVALID"
		messageKey = "InvalidArgumentError"
		metadatas = []map[string]string{
			{
				"messageKey": messageKey,
				"field":      "email",
			},
		}
	} else if {
          ...
        }

        // when adding multiple details
	for _, md := range anoterDetailData {
		for k, v := range md {
			metadatas = append(metadatas, map[string]string{
				"messageKey": messageKey,
				k:            v,
			})
		}
	}

	for _, md := range metadatas {
		details = append(details, &errdetails.ErrorInfo{
			Reason:   reason,
			Domain:   domain,
			Metadata: md,
		})
	}

	detailMessages := make([]protoiface.MessageV1, len(details))
	for i, d := range details {
		detailMessages[i] = d
	}
	dt, err := status.WithDetails(detailMessages...)
	if err != nil {
		return statusInternal.Err()
	}
	return dt.Err()
}
```

### Point2：Add a response using NewError() in case of error
If an error occurs, use `NewError()` and pass the status that occurred. 

If we want to add information, pass anoterDetailData as necessary.
```go
func validateGetAccountV2Request(req *accountproto.GetAccountV2Request, localizer locale.Localizer) error {
	if !verifyEmailFormat(req.Command.Email) {
		return NewError(statusInvalidEmail,
                                map[string]string{
                                  "value": req.Command.Email,
                                }
                       )
	}
...
```

### Point３：Remove localize massage code
Removed `LocalizedMessage` return code in case of error.

```go
func validateGetAccountV2Request(req *accountproto.GetAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(
			&errdetails.ErrorInfo{
				Reason: "INVALID",
				Domain: "account.bucketeer.io",
				Metadata: map[string]string{
					"messageKey": "InvalidArgumentError",
					"field": "email"
				},
			})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
...
```



## Frontend Updates

### 1. The temporary way to show the error message until releasing the all backend updates

Show the message in the `message` field when the error is returned from the backend.
Alghough it can show only English and the message is unclear, this is needed because the error format has been already updated in v2.1.1, and the right way might require more time to develop.

e.g. Error Response
```json
{
  "code": 2,
  "message": "rpc error: code = NotFound desc = account:account not found, account", // Use this
  "details": [ // Ignore for now
    {
      "@type": "type.googleapis.com/google.rpc.ErrorInfo",
      "reason": "UNKNOWN",
      "domain": "unknown.bucketeer.io",
      "metadata": {
        "messageKey": "unknown"
      }
    }
  ]
}
```

### 2. The complete way

Use the `messageKey`, other metadata, and message template files to show the complete message.
As the current implementation, react-i18next is useful to embed nouns.
We don't need to be aware of field names of the metadata except `messageKey` while developing the frontend.


## Development Steps

| Phase | Backend                                               | Frontend                                                                                                                                                       |
| ----- | ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 0.1   | Define the new error struct                           | -                                                                                                                                                              |
| 0.2   | Replace some of the existing errors to the new format | -                                                                                                                                                              |
| 1     | Start refining the error struct and replacing to it   | Update and release to show the `message` field temporarily. See [here](#1-the-temporary-way-to-show-the-error-message-until-releasing-the-all-backend-updates) |
| 2     | Complete replacing to the new struct                  | -                                                                                                                                                              |
| 3     | -                                                     | Update to show the complete message using the `messageKey`. See [here](#2-the-complete-way)                                                                    |
