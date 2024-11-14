# Storage challenges at API layer

API layer relies on MySQL for storage. Ideally, each layer is loosely coupled, and dependencies between each layer must be resolved.

The issues are as follows:

Issue1：RunTransaction, Begin, etc. are called and transaction processing is performed at the API layer.
Issue2：Assembling the Where clause and Order clause is done at the API layer.
Issue3：The types of DB instances held by each storage class are not uniform.

I propose the following solutions to the above issues.

## Solution for issue 1
RunTransaction, Begin, etc. are called and transaction processing is performed at the API layer.

### Example of code that is an issue
- `RunTransaction`
    - https://github.com/bucketeer-io/bucketeer/blob/main/pkg/account/api/account.go#L64
- `Begin`
    - https://github.com/bucketeer-io/bucketeer/blob/main/pkg/autoops/api/api.go#L185

### Solution
- Define the necessary API for each storage layer using Interface and abstract the storage layer.
    - In the storage layer, embed the interface for each storage and move the implementation necessary for DB operations to the storage layer.
    - The API layer calls the above Interface and does not directly perform DB operations.
- If transaction processing is required in the API layer, use RunTransactionAPI provided in the DB client to perform transaction processing.
    - If we want your API layer to handle transactions across different storages, we can do so by creating storage instances using the same DB client instance.
````
// Define the API required for storage with an interface
type AccountStorage interface {
	CreateAccount(ctx context.Context, a *domain.AccountV2) error
	UpdateAccount(ctx context.Context, a *domain.AccountV2) error
	DeleteAccount(ctx context.Context, a *domain.AccountV2) error
	GetAccount(ctx context.Context, email, organizationID string) (*domain.AccountV2, error)
  ...
}

// The API service maintains instances based on the required storage interface.
type AccountService struct {
	accountStorage    v2.AccountStorage     // storage instance
	publisher         publisher.Publisher
	opts              *options
	logger            *zap.Logger
}

func (s *AccountService) updateAccount(
	ctx context.Context,
	editor *eventproto.Editor,
	email, organizationID string,
) error {
  // No DB operations such as transaction processing are performed, just call Update.
	return s.accountStorage.UpdateAccount(ctx, account)
}

type Client interface {
  // Use when transaction processing is required
	RunInTransaction(ctx context.Context, tx Transaction, f func() error) error
}

func (c *client) RunInTransaction(ctx context.Context, f func() error) error {
	tx, err := c.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("client: begin tx: %w", err)
	}
	ctx = context.WithValue(ctx, transactionKey, tx)
	return c.runInTransaction(ctx, tx, f)
}

func (c *client) runInTransaction(ctx context.Context, tx Transaction, f func() error) error {
	var err error
	defer record()(operationRunInTransaction, &err)
	defer func() {
		if err != nil {
			tx.Rollback() // nolint:errcheck
		}
	}()
	if err = f(); err == nil {
		err = tx.Commit()
	}
	return err
}

// Called when executing a query on the storage layer
func (c *client) QueryExecer(ctx context.Context) mysql.QueryExecer {
	tx, ok := ctx.Value(transactionKey).(mysql.Transaction)
	if ok {
		return tx
	}
	return c
}

````
#### Transaction flow between different storages
<img width="1119" alt="_9___API層_MySQL依存解消" src="https://github.com/user-attachments/assets/1612123b-e484-453f-8eaf-8a3208588a7e">

## Solution for issue 2
Assembling the Where clause and Order clause is done at the API layer.

### Example of code that is an issue
- `AutoOps`
https://github.com/bucketeer-io/bucketeer/blob/main/pkg/autoops/api/api.go#L1115

### Solution
- Define a common condition specification structure(ListOptions) that specifies acquisition conditions in the List system, and specify conditions from the API layer via it.
 The storage layer converts it into conditions suitable for DB operations and uses it.

 ````
 // Options for specifying conditions structure
 type ListOptions struct {
 	Limit   int
 	Filters []ListFilter
 	Orders  []Order
 	Cursor  string
 }

 type ListFilter struct {
 	Field    string
 	Operator Operator
 	Value    interface{}
 }

 type Operator int

 const (
 	OperatorEqual = 1
 	OperatorNotEqual
 	....
 )

 type Order struct {
 	Field     string
 	Direction OrderDirection
 }

 type OrderDirection int

 const (
 	Asc OrderDirection = 1
 	Desc
 )


 // Example
 ListOptions{
 	Limit: 10,
 	Filters: []ListFilter{
 		{
 			Field:    "organization_id",
 			Operator: OperatorEqual,
 			Value:    "organization_id-1",
 		},
 	},
 	Orders: []Order{
 		{
 			Field: "email",
 			Direction: Asc,
 		},
 	},
 }

 // Example for List
 type AccountStorage interface {
	ListAccounts(
		ctx context.Context,
		listOptions *storage.ListOptions,   // Condition specification options
	) ([]*proto.AccountV2, int, int64, error)
}

func (s *accountStorage) ListAccountsMySQL(
	ctx context.Context,
	listOptions *storage.ListOptions,
) ([]*proto.AccountV2, int, int64, error) {
  whereParts := s.getWhereParts(listOptions.filters)   // Convert to the Where clause
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)

	orderBySQL := mysql.ConstructOrderBySQLString(listOptions.orders)  //Convert to the Order clause
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(listOptions.limit, listOptions.cursor)
	query := fmt.Sprintf(
		selectAccountsV2SQL,
		whereSQL,
		orderBySQL,
		limitOffsetSQL,
	)
	rows, err := s.qe(ctx).QueryContext(ctx, query, whereArgs...)

・・・・
}
 ````

 ## Solution for issue 3
The types of DB instances held by each storage class are not uniform.

### Example of code that is an issue
- autoOps holds QueryExecer
https://github.com/bucketeer-io/bucketeer/blob/main/pkg/autoops/storage/v2/auto_ops_rule.go#L59
- account holds Client
https://github.com/bucketeer-io/bucketeer/blob/main/pkg/account/storage/v2/storage.go#L55

### Solution
- Each storage tier holds the DB client instances required for each storage.
  In the storage layer, DB operations are performed via the DB client.

  ````
  type accountStorage struct {
  	client mysql.Client  // MySQL Client
  }

  func (s *accountStorage) UpdateAccount(ctx context.Context, a *domain.AccountV2) error {
    result, err := s.client.QueryExecer(ctx).ExecContext(
                             ctx,
                             updateAccountV2SQL,
                             a.Name,
                             a.FirstName,
                             a.LastName,
                             a.Language,
                             a.AvatarImageUrl,
                             a.AvatarFileType,
                             a.AvatarImage,
                             int32(a.OrganizationRole),
                             mysql.JSONObject{Val: a.EnvironmentRoles},
                             a.Disabled,
                             a.UpdatedAt,
                             a.LastSeen,
                             mysql.JSONObject{Val: a.SearchFilters},
                             a.Email,
                             a.OrganizationId,
                             )
    if err != nil {
      return err
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
      return err
    }
    if rowsAffected != 1 {
      return ErrAccountUnexpectedAffectedRows
    }
    return nil
  }
  ````

## Associated storage
- [ ]  account
- [ ]  auditlog
- [ ]  autoops
- [ ]  environment
- [ ]  eventcounter
- [ ]  experiment
- [ ]  feature
- [ ]  notification
- [ ]  push
- [ ]  experimentcalculator
- [ ]  mau
- [ ]  opsevent
- [ ]  subscriber
