// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	notificationstorage "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

var (
	//go:embed sql/insert_notification.sql
	insertNotificationSQL string
	//go:embed sql/insert_notification_localization.sql
	insertNotificationLocalizationSQL string
	//go:embed sql/select_notification.sql
	selectNotificationSQL string
	//go:embed sql/update_notification.sql
	updateNotificationSQL string
	//go:embed sql/publish_notification.sql
	publishNotificationSQL string
	//go:embed sql/delete_notification_localizations.sql
	deleteNotificationLocalizationsSQL string
	//go:embed sql/delete_notification.sql
	deleteNotificationSQL string
	//go:embed sql/select_draft_notifications.sql
	selectDraftNotificationsSQL string
	//go:embed sql/count_draft_notifications.sql
	countDraftNotificationsSQL string
	//go:embed sql/select_notification_localizations.sql
	selectNotificationLocalizationsSQL string
	//go:embed sql/select_notifications.sql
	selectNotificationsSQL string
	//go:embed sql/count_notifications.sql
	countNotificationsSQL string
	//go:embed sql/select_read_notification_ids.sql
	selectReadNotificationIDsSQL string
	//go:embed sql/select_earliest_account_created_at.sql
	selectEarliestAccountCreatedAtSQL string
)

// readNotificationExistsSubquery correlates a viewer's read marker with the
// outer notification row for EXISTS / NOT EXISTS read-status filters; the
// $%d template is bound at query-construction time.
const readNotificationExistsSubquery = "SELECT 1 FROM notification_read " +
	"WHERE notification_read.notification_id = notification.id AND notification_read.email = $%d"

type notificationStorage struct {
	qe pgstorage.QueryExecer
}

func NewNotificationStorage(qe pgstorage.QueryExecer) notificationstorage.NotificationStorage {
	return &notificationStorage{qe: qe}
}

func (s *notificationStorage) CreateAdminNotification(
	ctx context.Context,
	notification *domain.Notification,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertNotificationSQL,
		notification.Id,
		int32(notification.Status),
		notification.CreatedBy,
		notification.LastEditedBy,
		notification.CreatedAt,
		notification.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return notificationstorage.ErrNotificationAlreadyExists
		}
		return err
	}
	for _, l := range notification.Localizations {
		_, err := s.qe.ExecContext(
			ctx,
			insertNotificationLocalizationSQL,
			notification.Id,
			l.Language,
			pgstorage.JSONObject{Val: l.Tags},
			l.Title,
			l.Content,
		)
		if err != nil {
			if errors.Is(err, pgstorage.ErrDuplicateEntry) {
				return notificationstorage.ErrNotificationAlreadyExists
			}
			return err
		}
	}
	return nil
}

func (s *notificationStorage) GetAdminNotification(
	ctx context.Context,
	id string,
) (*domain.Notification, error) {
	notification := proto.Notification{}
	var status int32
	err := s.qe.QueryRowContext(
		ctx,
		selectNotificationSQL,
		id,
	).Scan(
		&notification.Id,
		&status,
		&notification.CreatedBy,
		&notification.LastEditedBy,
		&notification.PublishedBy,
		&notification.PublishedAt,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, notificationstorage.ErrNotificationNotFound
		}
		return nil, err
	}
	notification.Status = proto.Notification_Status(status)
	if err := s.fillLocalizations(ctx, []*proto.Notification{&notification}); err != nil {
		return nil, err
	}
	return &domain.Notification{Notification: &notification}, nil
}

func (s *notificationStorage) UpdateAdminNotification(
	ctx context.Context,
	notification *domain.Notification,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateNotificationSQL,
		notification.LastEditedBy,
		notification.UpdatedAt,
		notification.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return notificationstorage.ErrNotificationNotFound
	}
	if _, err := s.qe.ExecContext(
		ctx,
		deleteNotificationLocalizationsSQL,
		notification.Id,
	); err != nil {
		return err
	}
	for _, l := range notification.Localizations {
		_, err := s.qe.ExecContext(
			ctx,
			insertNotificationLocalizationSQL,
			notification.Id,
			l.Language,
			pgstorage.JSONObject{Val: l.Tags},
			l.Title,
			l.Content,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *notificationStorage) PublishAdminNotification(
	ctx context.Context,
	notification *domain.Notification,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		publishNotificationSQL,
		int32(notification.Status),
		notification.PublishedBy,
		notification.PublishedAt,
		notification.UpdatedAt,
		notification.Id,
		int32(proto.Notification_DRAFT),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// The row is no longer a live draft: a concurrent request
		// published (or soft-deleted) it after our read.
		return notificationstorage.ErrNotificationAlreadyPublished
	}
	return nil
}

func (s *notificationStorage) DeleteAdminNotification(
	ctx context.Context,
	id, lastEditedBy string,
	updatedAt int64,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteNotificationSQL,
		lastEditedBy,
		updatedAt,
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return notificationstorage.ErrNotificationNotFound
	}
	return nil
}

func listDraftAdminNotificationsOrders(
	orderBy proto.ListDraftAdminNotificationsRequest_OrderBy,
	orderDirection proto.ListDraftAdminNotificationsRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListDraftAdminNotificationsRequest_DEFAULT,
		proto.ListDraftAdminNotificationsRequest_CREATED_AT:
		column = "notification.created_at"
	case proto.ListDraftAdminNotificationsRequest_UPDATED_AT:
		column = "notification.updated_at"
	default:
		return nil, notificationstorage.ErrInvalidListDraftAdminNotificationsOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListDraftAdminNotificationsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}

func listDraftAdminNotificationsFilters(
	searchKeyword string,
) ([]*pgstorage.Filter, *pgstorage.SearchQuery) {
	filters := []*pgstorage.Filter{
		{
			Column:   "notification.status",
			Operator: pgstorage.OperatorEqual,
			Value:    int32(proto.Notification_DRAFT),
		},
		{
			Column:   "notification.deleted",
			Operator: pgstorage.OperatorEqual,
			Value:    false,
		},
	}
	var searchQuery *pgstorage.SearchQuery
	if searchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{
				"notification_localization.title",
				"notification_localization.content",
			},
			Keyword: searchKeyword,
		}
	}
	return filters, searchQuery
}

func (s *notificationStorage) ListDraftAdminNotifications(
	ctx context.Context,
	p notificationstorage.ListDraftAdminNotificationsParams,
) ([]*proto.Notification, int, int64, error) {
	orders, err := listDraftAdminNotificationsOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, 0, 0, err
	}
	filters, searchQuery := listDraftAdminNotificationsFilters(p.SearchKeyword)
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, 0, 0, notificationstorage.ErrInvalidListDraftAdminNotificationsCursor
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	options := &pgstorage.ListOptions{
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       limit,
		Offset:      offset,
	}
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(options.CreateWhereParts())
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf(selectDraftNotificationsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	notifications := make([]*proto.Notification, 0, limit)
	for rows.Next() {
		notification := proto.Notification{}
		var status int32
		err := rows.Scan(
			&notification.Id,
			&status,
			&notification.CreatedBy,
			&notification.LastEditedBy,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		notification.Status = proto.Notification_Status(status)
		notifications = append(notifications, &notification)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	if err := s.fillLocalizations(ctx, notifications); err != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(notifications)
	var totalCount int64
	countQuery := fmt.Sprintf(countDraftNotificationsSQL, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return notifications, nextOffset, totalCount, nil
}

func listNotificationsOrders(
	orderBy proto.ListNotificationsRequest_OrderBy,
	orderDirection proto.ListNotificationsRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListNotificationsRequest_DEFAULT,
		proto.ListNotificationsRequest_PUBLISHED_AT:
		column = "notification.published_at"
	default:
		return nil, notificationstorage.ErrInvalidListNotificationsOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListNotificationsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}

func (s *notificationStorage) ListNotifications(
	ctx context.Context,
	p notificationstorage.ListNotificationsParams,
) ([]*proto.Notification, int, int64, error) {
	orders, err := listNotificationsOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, 0, 0, err
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, 0, 0, notificationstorage.ErrInvalidListNotificationsCursor
	}
	filters := []*pgstorage.Filter{
		{
			Column:   "notification.status",
			Operator: pgstorage.OperatorEqual,
			Value:    int32(proto.Notification_PUBLISHED),
		},
		{
			Column:   "notification.deleted",
			Operator: pgstorage.OperatorEqual,
			Value:    false,
		},
	}
	if p.PublishedAtFrom > 0 {
		filters = append(filters, &pgstorage.Filter{
			Column:   "notification.published_at",
			Operator: pgstorage.OperatorGreaterThanOrEqual,
			Value:    p.PublishedAtFrom,
		})
	}
	if p.PublishedAtTo > 0 {
		filters = append(filters, &pgstorage.Filter{
			Column:   "notification.published_at",
			Operator: pgstorage.OperatorLessThanOrEqual,
			Value:    p.PublishedAtTo,
		})
	}
	var existsFilters []*pgstorage.ExistsFilter
	switch p.ReadStatus {
	case proto.ListNotificationsRequest_UNREAD:
		// New users do not inherit history: unread only considers
		// notifications published after the account was created.
		accountCreatedAt, err := s.earliestAccountCreatedAt(ctx, p.Email)
		if err != nil {
			return nil, 0, 0, err
		}
		if accountCreatedAt > 0 {
			filters = append(filters, &pgstorage.Filter{
				Column:   "notification.published_at",
				Operator: pgstorage.OperatorGreaterThanOrEqual,
				Value:    accountCreatedAt,
			})
		}
		existsFilters = append(existsFilters, &pgstorage.ExistsFilter{
			Subquery:  readNotificationExistsSubquery,
			NotExists: true,
			Values:    []interface{}{p.Email},
		})
	case proto.ListNotificationsRequest_READ:
		existsFilters = append(existsFilters, &pgstorage.ExistsFilter{
			Subquery: readNotificationExistsSubquery,
			Values:   []interface{}{p.Email},
		})
	}
	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{
				"notification_localization.title",
				"notification_localization.content",
			},
			Keyword: p.SearchKeyword,
		}
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	options := &pgstorage.ListOptions{
		Filters:       filters,
		ExistsFilters: existsFilters,
		SearchQuery:   searchQuery,
		Orders:        orders,
		Limit:         limit,
		Offset:        offset,
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectNotificationsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	notifications := make([]*proto.Notification, 0, limit)
	for rows.Next() {
		notification := proto.Notification{}
		var status int32
		err := rows.Scan(
			&notification.Id,
			&status,
			&notification.CreatedBy,
			&notification.LastEditedBy,
			&notification.PublishedBy,
			&notification.PublishedAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		notification.Status = proto.Notification_Status(status)
		notifications = append(notifications, &notification)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	if err := s.fillLocalizations(ctx, notifications); err != nil {
		return nil, 0, 0, err
	}
	readIDs, err := s.readNotificationIDs(ctx, p.Email, notifications)
	if err != nil {
		return nil, 0, 0, err
	}
	for _, n := range notifications {
		n.Localization = domain.ResolveLocalization(n.Localizations, p.Language)
		n.Localizations = nil
		_, n.Read = readIDs[n.Id]
	}
	nextOffset := offset + len(notifications)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(countNotificationsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return notifications, nextOffset, totalCount, nil
}

// IsNotificationRead reports whether the viewer has read the notification.
func (s *notificationStorage) IsNotificationRead(
	ctx context.Context,
	id, email string,
) (bool, error) {
	readIDs, err := s.readNotificationIDs(ctx, email, []*proto.Notification{{Id: id}})
	if err != nil {
		return false, err
	}
	_, ok := readIDs[id]
	return ok, nil
}

// readNotificationIDs returns which of the given notifications the viewer
// has read; the query is bounded by the page size.
func (s *notificationStorage) readNotificationIDs(
	ctx context.Context,
	email string,
	notifications []*proto.Notification,
) (map[string]struct{}, error) {
	if len(notifications) == 0 {
		return map[string]struct{}{}, nil
	}
	ids := make([]interface{}, 0, len(notifications))
	for _, n := range notifications {
		ids = append(ids, n.Id)
	}
	options := &pgstorage.ListOptions{
		Filters: []*pgstorage.Filter{
			{
				Column:   "notification_read.email",
				Operator: pgstorage.OperatorEqual,
				Value:    email,
			},
		},
		InFilters: []*pgstorage.InFilter{
			{
				Column: "notification_read.notification_id",
				Values: ids,
			},
		},
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectReadNotificationIDsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	readIDs := map[string]struct{}{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		readIDs[id] = struct{}{}
	}
	return readIDs, rows.Err()
}

func (s *notificationStorage) earliestAccountCreatedAt(
	ctx context.Context,
	email string,
) (int64, error) {
	var createdAt int64
	err := s.qe.QueryRowContext(
		ctx,
		selectEarliestAccountCreatedAtSQL,
		email,
	).Scan(&createdAt)
	if err != nil {
		return 0, err
	}
	return createdAt, nil
}

func (s *notificationStorage) fillLocalizations(
	ctx context.Context,
	notifications []*proto.Notification,
) error {
	if len(notifications) == 0 {
		return nil
	}
	ids := make([]interface{}, 0, len(notifications))
	byID := make(map[string]*proto.Notification, len(notifications))
	for _, n := range notifications {
		ids = append(ids, n.Id)
		byID[n.Id] = n
	}
	options := &pgstorage.ListOptions{
		InFilters: []*pgstorage.InFilter{
			{
				Column: "notification_localization.notification_id",
				Values: ids,
			},
		},
		Orders: []*pgstorage.Order{
			pgstorage.NewOrder("notification_localization.language", pgstorage.OrderDirectionAsc),
		},
	}
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(options.CreateWhereParts())
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	query := fmt.Sprintf(selectNotificationLocalizationsSQL, whereSQL, orderBySQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var notificationID string
		localization := proto.NotificationLocalization{}
		err := rows.Scan(
			&notificationID,
			&localization.Language,
			&pgstorage.JSONObject{Val: &localization.Tags},
			&localization.Title,
			&localization.Content,
		)
		if err != nil {
			return err
		}
		if n, ok := byID[notificationID]; ok {
			n.Localizations = append(n.Localizations, &localization)
		}
	}
	return rows.Err()
}
