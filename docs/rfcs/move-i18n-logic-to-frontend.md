# Summary

Now, we have part of the i18n implementation on both sides, front and backend. 

Because the i18n implementation in all APIs made the backend code too hard to read, we decided to do this only on the front end.

Note that this RFC only describes content related to the backend.

[Issue](https://github.com/bucketeer-io/bucketeer/issues/1253)


## Response design
Currently, I am using GRPC's ErrorDetail to return a localized message based on LocalizedMessage.

[GRPC's LocalizedMessage](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto#L290)

Ex: Validation error when GetAccount

https://github.com/bucketeer-io/bucketeer/blob/main/pkg/account/api/validation.go#L624

Design to return by utilizing ErrorInfo defined in GRPC's ErrorDetail.

[GRPC's ErrorInfo]https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto#L51

If multiple errors are returned, multiple ErrorInfos are returned.

Response Ex：
```json
{
  "reason": "INVALID"
  "domain": "account.bucketeer.io",
  "metadata": {
      "messageKey": "account.invalid.format",
      "field": "email",
      "value": "email.com",
  }
}
```
| Key | Explanation | Example |
|:---|:---|:---|
|reason|The reason of the error.|"reason": "INVALID"|
|domain|The error domain is typically the registered service name of the tool or product that generates the error. |"domain": "account.bucketeer.io"  // account package|
|metadata| Additional structured details about this error. |- |
|messageKey| Key to identify message content.<br> Format: ［error package name］.［error type］.(error characteristics)<br>※ Grant error characteristics only when necessary. |ex1) "messageKey": "account.invalid" // Invalid error in account package<br> ex2) "messageKey": "account.invalid.empty" // Invalid error with empty value in account package |
|field| Send field information in a message. Granted only when needed. |"field": "email" |
|value| Send value information in a message. Granted only when needed. |"value": "email.com" |


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
		messageKey = "account.invalid.empty"
		metadatas = []map[string]string{
			{
				"messageKey": messageKey,
				"feild":      "email",
			},
		}
	} else if status == statusInvalidEmail {
		reason = "INVALID"
		messageKey = "account.invalid.format"
		metadatas = []map[string]string{
			{
				"messageKey": messageKey,
				"feild":      "email",
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
					"messageKey": "account.invalid.empty",
					"feild":      "email",
				},
			})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
...
```

## Release Steps
1. Releases a process that returns an ErrorInfo for each package for The backend.
2. Supports multilingual ization based on the ErrorInfo process in step 1 with Frontend.

※ Currently, the front localization information acquisition process seems to be implemented using the code below, so it may be a good idea to focus on modifying that area.
https://github.com/bucketeer-io/bucketeer/blob/main/ui/web-v2/src/grpc/messages.ts

3. Remove the LocalizedMessage code in backend for the error handled in step 2.
