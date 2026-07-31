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

package mysql

import (
	"context"
	_ "embed"
	"errors"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	notificationstorage "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
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
// outer notification row for EXISTS / NOT EXISTS read-status filters.
const readNotificationExistsSubquery = "SELECT 1 FROM notification_read " +
	"WHERE notification_read.notification_id = notification.id AND notification_read.email = ?"

type notificationStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewNotificationStorage(qe mysqlstorage.QueryExecer) notificationstorage.NotificationStorage {
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
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
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
			mysqlstorage.JSONObject{Val: l.Tags},
			l.Title,
			l.Content,
		)
		if err != nil {
			if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
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
		if errors.Is(err, mysqlstorage.ErrNoRows) {
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
			mysqlstorage.JSONObject{Val: l.Tags},
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
) ([]*mysqlstorage.Order, error) {
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
	direction := mysqlstorage.OrderDirectionAsc
	if orderDirection == proto.ListDraftAdminNotificationsRequest_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}
	return []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)}, nil
}

func listDraftAdminNotificationsFilters(
	searchKeyword string,
) ([]*mysqlstorage.FilterV2, *mysqlstorage.SearchQuery) {
	filters := []*mysqlstorage.FilterV2{
		{
			Column:   "notification.status",
			Operator: mysqlstorage.OperatorEqual,
			Value:    int32(proto.Notification_DRAFT),
		},
		{
			Column:   "notification.deleted",
			Operator: mysqlstorage.OperatorEqual,
			Value:    false,
		},
	}
	var searchQuery *mysqlstorage.SearchQuery
	if searchKeyword != "" {
		searchQuery = &mysqlstorage.SearchQuery{
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
	options := &mysqlstorage.ListOptions{
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       limit,
		Offset:      offset,
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectDraftNotificationsSQL, options)
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
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countDraftNotificationsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return notifications, nextOffset, totalCount, nil
}

func listNotificationsOrders(
	orderBy proto.ListNotificationsRequest_OrderBy,
	orderDirection proto.ListNotificationsRequest_OrderDirection,
) ([]*mysqlstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListNotificationsRequest_DEFAULT,
		proto.ListNotificationsRequest_PUBLISHED_AT:
		column = "notification.published_at"
	default:
		return nil, notificationstorage.ErrInvalidListNotificationsOrderBy
	}
	direction := mysqlstorage.OrderDirectionAsc
	if orderDirection == proto.ListNotificationsRequest_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}
	return []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)}, nil
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
	filters := []*mysqlstorage.FilterV2{
		{
			Column:   "notification.status",
			Operator: mysqlstorage.OperatorEqual,
			Value:    int32(proto.Notification_PUBLISHED),
		},
		{
			Column:   "notification.deleted",
			Operator: mysqlstorage.OperatorEqual,
			Value:    false,
		},
	}
	if p.PublishedAtFrom > 0 {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "notification.published_at",
			Operator: mysqlstorage.OperatorGreaterThanOrEqual,
			Value:    p.PublishedAtFrom,
		})
	}
	if p.PublishedAtTo > 0 {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "notification.published_at",
			Operator: mysqlstorage.OperatorLessThanOrEqual,
			Value:    p.PublishedAtTo,
		})
	}
	var existsFilters []*mysqlstorage.ExistsFilter
	switch p.ReadStatus {
	case proto.ListNotificationsRequest_UNREAD:
		// New users do not inherit history: unread only considers
		// notifications published after the account was created.
		accountCreatedAt, err := s.earliestAccountCreatedAt(ctx, p.Email)
		if err != nil {
			return nil, 0, 0, err
		}
		if accountCreatedAt > 0 {
			filters = append(filters, &mysqlstorage.FilterV2{
				Column:   "notification.published_at",
				Operator: mysqlstorage.OperatorGreaterThanOrEqual,
				Value:    accountCreatedAt,
			})
		}
		existsFilters = append(existsFilters, &mysqlstorage.ExistsFilter{
			Subquery:  readNotificationExistsSubquery,
			NotExists: true,
			Values:    []interface{}{p.Email},
		})
	case proto.ListNotificationsRequest_READ:
		existsFilters = append(existsFilters, &mysqlstorage.ExistsFilter{
			Subquery: readNotificationExistsSubquery,
			Values:   []interface{}{p.Email},
		})
	}
	var searchQuery *mysqlstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &mysqlstorage.SearchQuery{
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
	options := &mysqlstorage.ListOptions{
		Filters:       filters,
		ExistsFilters: existsFilters,
		SearchQuery:   searchQuery,
		Orders:        orders,
		Limit:         limit,
		Offset:        offset,
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectNotificationsSQL, options)
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
	readIDs, err := s.readNotificationIDs(ctx, p.Email)
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
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countNotificationsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return notifications, nextOffset, totalCount, nil
}

// readNotificationIDs returns the set of notification ids the viewer has read.
func (s *notificationStorage) readNotificationIDs(
	ctx context.Context,
	email string,
) (map[string]struct{}, error) {
	options := &mysqlstorage.ListOptions{
		Filters: []*mysqlstorage.FilterV2{
			{
				Column:   "notification_read.email",
				Operator: mysqlstorage.OperatorEqual,
				Value:    email,
			},
		},
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectReadNotificationIDsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := map[string]struct{}{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = struct{}{}
	}
	return ids, rows.Err()
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
	options := &mysqlstorage.ListOptions{
		InFilters: []*mysqlstorage.InFilter{
			{
				Column: "notification_localization.notification_id",
				Values: ids,
			},
		},
		Orders: []*mysqlstorage.Order{
			mysqlstorage.NewOrder("notification_localization.language", mysqlstorage.OrderDirectionAsc),
		},
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectNotificationLocalizationsSQL, options)
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
			&mysqlstorage.JSONObject{Val: &localization.Tags},
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
