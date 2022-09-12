/*
 *    Copyright 2022 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package sqltrace

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/chenquan/sqlplus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	_                         sqlplus.Hook = (*Hook)(nil)
	sqlMethodAttributeKey                  = attribute.Key("sql.method")
	sqlAttributeKey                        = attribute.Key("sql")
	sqlArgsLengthAttributeKey              = attribute.Key("sql.args_length")
)

type (
	Hook struct {
		c Config
	}
)

const (
	spanName = "sql"
)

func NewTraceHook(c Config) *Hook {
	StartAgent(c)
	return &Hook{c: c}
}

func (h *Hook) BeforeConnect(ctx context.Context, err error) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "connect")

	return ctx, err
}

func (h *Hook) AfterConnect(ctx context.Context, dc driver.Conn, err error) (context.Context, driver.Conn, error) {
	h.endSpan(ctx, err)

	return ctx, dc, err
}

func (h *Hook) BeforeExecContext(ctx context.Context, query string, args []driver.NamedValue, err error) (context.Context, string, []driver.NamedValue, error) {
	ctx, span := h.startSpan(ctx, "exec")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	return ctx, query, args, err
}

func (h *Hook) AfterExecContext(ctx context.Context, _ string, _ []driver.NamedValue, r driver.Result, err error) (context.Context, driver.Result, error) {
	h.endSpan(ctx, err)

	return ctx, r, err
}

func (h *Hook) BeforeBeginTx(ctx context.Context, opts driver.TxOptions, err error) (context.Context, driver.TxOptions, error) {
	ctx, _ = h.startSpan(ctx, "begin")

	return ctx, opts, err
}

func (h *Hook) AfterBeginTx(ctx context.Context, _ driver.TxOptions, dd driver.Tx, err error) (context.Context, driver.Tx, error) {
	h.endSpan(ctx, err)

	return ctx, dd, err
}

func (h *Hook) BeforeQueryContext(ctx context.Context, query string, args []driver.NamedValue, err error) (context.Context, string, []driver.NamedValue, error) {
	ctx, span := h.startSpan(ctx, "query")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	return ctx, query, args, err
}

func (h *Hook) AfterQueryContext(ctx context.Context, _ string, _ []driver.NamedValue, rows driver.Rows, err error) (context.Context, driver.Rows, error) {
	h.endSpan(ctx, err)

	return ctx, rows, err
}

func (h *Hook) BeforePrepareContext(ctx context.Context, query string, err error) (context.Context, string, error) {
	ctx, span := h.startSpan(ctx, "prepare")
	span.SetAttributes(sqlAttributeKey.String(query))

	return ctx, query, err
}

func (h *Hook) AfterPrepareContext(ctx context.Context, _ string, s driver.Stmt, err error) (context.Context, driver.Stmt, error) {
	h.endSpan(ctx, err)

	return ctx, s, err
}

func (h *Hook) BeforeCommit(ctx context.Context, err error) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "commit")

	return ctx, err
}

func (h *Hook) AfterCommit(ctx context.Context, err error) (context.Context, error) {
	h.endSpan(ctx, err)

	return ctx, err
}

func (h *Hook) BeforeRollback(ctx context.Context, err error) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "rollback")

	return ctx, err
}

func (h *Hook) AfterRollback(ctx context.Context, err error) (context.Context, error) {
	h.endSpan(ctx, err)

	return ctx, err
}

func (h *Hook) BeforeStmtQueryContext(ctx context.Context, query string, args []driver.NamedValue, err error) (context.Context, []driver.NamedValue, error) {
	ctx, span := h.startSpan(ctx, "stmtQuery")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	return ctx, args, err
}

func (h *Hook) AfterStmtQueryContext(ctx context.Context, _ string, _ []driver.NamedValue, rows driver.Rows, err error) (context.Context, driver.Rows, error) {
	h.endSpan(ctx, err)

	return ctx, rows, err
}

func (h *Hook) BeforeStmtExecContext(ctx context.Context, query string, args []driver.NamedValue, err error) (context.Context, []driver.NamedValue, error) {
	ctx, span := h.startSpan(ctx, "stmtExec")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	return ctx, args, err
}

func (h *Hook) AfterStmtExecContext(ctx context.Context, _ string, _ []driver.NamedValue, r driver.Result, err error) (context.Context, driver.Result, error) {
	h.endSpan(ctx, err)

	return ctx, r, err
}

func (h *Hook) startSpan(ctx context.Context, method string) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer(h.c.Name)
	spanStartOptions := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(semconv.DBSystemKey.String(h.c.DataSourceName)),
	}

	prepareContext := sqlplus.PrepareContextFromContext(ctx)
	if prepareContext != nil {
		spanStartOptions = append(spanStartOptions, trace.WithLinks(trace.LinkFromContext(prepareContext)))
	}

	txContext := sqlplus.TxContextFromContext(ctx)
	if txContext != nil {
		spanStartOptions = append(spanStartOptions, trace.WithLinks(trace.LinkFromContext(txContext)))
	}

	start, span := tracer.Start(ctx,
		spanName,
		spanStartOptions...,
	)
	span.SetAttributes(sqlMethodAttributeKey.String(method))

	return start, span
}

func (h *Hook) endSpan(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if err == nil || err == sql.ErrNoRows {
		span.SetStatus(codes.Ok, "")
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
