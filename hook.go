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
		hook sqlplus.Hook
		c    Config
	}
)

const (
	spanName = "sql"
)

func (h *Hook) SetHook(hook sqlplus.Hook) {
	h.hook = hook
}

func (h *Hook) BeforeConnect(ctx context.Context) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "connect")

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeConnect(ctx)
	}

	return ctx, err
}

func (h *Hook) AfterConnect(ctx context.Context, dc driver.Conn, err error) (context.Context, driver.Conn, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, dc, err = h.hook.AfterConnect(ctx, dc, err)
	}

	return ctx, dc, err
}

func (h *Hook) BeforeExecContext(ctx context.Context, query string, args []driver.NamedValue) (context.Context, error) {
	ctx, span := h.startSpan(ctx, "exec")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeExecContext(ctx, query, args)
	}

	return ctx, err
}

func (h *Hook) AfterExecContext(ctx context.Context, query string, args []driver.NamedValue, r driver.Result, err error) (context.Context, driver.Result, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, r, err = h.hook.AfterExecContext(ctx, query, args, r, err)
	}

	return ctx, r, err
}

func (h *Hook) BeforeBeginTx(ctx context.Context, opts driver.TxOptions) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "begin")

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeBeginTx(ctx, opts)
	}

	return ctx, err
}

func (h *Hook) AfterBeginTx(ctx context.Context, opts driver.TxOptions, dd driver.Tx, err error) (context.Context, driver.Tx, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, dd, err = h.hook.AfterBeginTx(ctx, opts, dd, err)
	}

	return ctx, dd, err
}

func (h *Hook) BeforeQueryContext(ctx context.Context, query string, args []driver.NamedValue) (context.Context, error) {
	ctx, span := h.startSpan(ctx, "query")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeQueryContext(ctx, query, args)
	}

	return ctx, err
}

func (h *Hook) AfterQueryContext(ctx context.Context, query string, args []driver.NamedValue, rows driver.Rows, err error) (context.Context, driver.Rows, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, rows, err = h.hook.AfterQueryContext(ctx, query, args, rows, err)
	}

	return ctx, rows, err
}

func (h *Hook) BeforePrepareContext(ctx context.Context, query string) (context.Context, error) {
	ctx, span := h.startSpan(ctx, "prepare")
	span.SetAttributes(sqlAttributeKey.String(query))

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforePrepareContext(ctx, query)
	}

	return ctx, err
}

func (h *Hook) AfterPrepareContext(ctx context.Context, query string, s driver.Stmt, err error) (context.Context, driver.Stmt, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, s, err = h.hook.AfterPrepareContext(ctx, query, s, err)
	}

	return ctx, s, err
}

func (h *Hook) BeforeCommit(ctx context.Context) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "commit")

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeCommit(ctx)
	}

	return ctx, err
}

func (h *Hook) AfterCommit(ctx context.Context, err error) (context.Context, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, err = h.hook.AfterCommit(ctx, err)
	}

	return ctx, err
}

func (h *Hook) BeforeRollback(ctx context.Context) (context.Context, error) {
	ctx, _ = h.startSpan(ctx, "rollback")

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeRollback(ctx)
	}

	return ctx, err
}

func (h *Hook) AfterRollback(ctx context.Context, err error) (context.Context, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, err = h.hook.AfterRollback(ctx, err)
	}

	return ctx, err
}

func (h *Hook) BeforeStmtQueryContext(ctx context.Context, query string, args []driver.NamedValue) (context.Context, error) {
	ctx, span := h.startSpan(ctx, "stmtQuery")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeStmtQueryContext(ctx, query, args)
	}

	return ctx, err
}

func (h *Hook) AfterStmtQueryContext(ctx context.Context, query string, args []driver.NamedValue, rows driver.Rows, err error) (context.Context, driver.Rows, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, rows, err = h.hook.AfterStmtQueryContext(ctx, query, args, rows, err)
	}

	return ctx, rows, err
}

func (h *Hook) BeforeStmtExecContext(ctx context.Context, query string, args []driver.NamedValue) (context.Context, error) {
	ctx, span := h.startSpan(ctx, "stmtExec")
	span.SetAttributes(sqlAttributeKey.String(query))
	span.SetAttributes(sqlArgsLengthAttributeKey.Int(len(args)))

	var err error
	if h.hook != nil {
		ctx, err = h.hook.BeforeStmtExecContext(ctx, query, args)
	}

	return ctx, err
}

func (h *Hook) AfterStmtExecContext(ctx context.Context, query string, args []driver.NamedValue, r driver.Result, err error) (context.Context, driver.Result, error) {
	h.endSpan(ctx, err)

	if h.hook != nil {
		ctx, r, err = h.hook.AfterStmtExecContext(ctx, query, args, r, err)
	}

	return ctx, r, err
}

func (h *Hook) startSpan(ctx context.Context, method string) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer(h.c.Name)
	start, span := tracer.Start(ctx,
		spanName,
		trace.WithSpanKind(trace.SpanKindClient),
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
