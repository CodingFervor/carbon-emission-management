package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/CodingFervor/carbon-emission-management/pkg/logger"
)

// Page holds the parameters of a paged query.
type Page struct {
	Page     int // 1-based page number
	PageSize int // items per page
}

// Normalize clamps the paging values to safe bounds and returns the SQL OFFSET/LIMIT offset.
func (p *Page) Normalize() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 200 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

// PageFromValues builds a Page from raw query-string values.
func PageFromValues(page, pageSize int) Page {
	return Page{Page: page, PageSize: pageSize}
}

// DB exposes the shared database handle to all repositories.
type DB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *DB { return &DB{db: db} }

func (d *DB) Sql() *sql.DB { return d.db }

// GenericRepo provides common CRUD operations for entities backed by a table
// whose columns map 1:1 to struct fields via the `db` tag. The dest factory
// returns a pointer to a fresh struct of type T for row scanning.
type GenericRepo[T any] struct {
	db       *sql.DB
	table    string
	columns  []string
	scanCols []string // columns without id, created_at, updated_at (for inserts)
	dest     func() *T
}

// NewGenericRepo builds a GenericRepo. It derives the column list from the `db`
// tags of the struct returned by dest (excluding private fields).
func NewGenericRepo[T any](db *sql.DB, table string, dest func() *T) *GenericRepo[T] {
	r := &GenericRepo[T]{db: db, table: table, dest: dest}
	r.columns = structColumns(dest())
	r.scanCols = r.columns
	return r
}

func (r *GenericRepo[T]) DB() *sql.DB   { return r.db }
func (r *GenericRepo[T]) Table() string { return r.table }

// Count returns the total number of rows matching an optional WHERE clause.
func (r *GenericRepo[T]) Count(where string, args ...any) (int64, error) {
	q := "SELECT COUNT(*) FROM " + r.table
	if where != "" {
		q += " WHERE " + where
	}
	var n int64
	if err := r.db.QueryRow(q, args...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

// List returns a paged slice of entities ordered by id desc.
func (r *GenericRepo[T]) List(p Page, where string, args ...any) ([]T, int64, error) {
	total, err := r.Count(where, args...)
	if err != nil {
		return nil, 0, err
	}
	offset := p.Normalize()
	q := "SELECT " + strings.Join(r.columns, ", ") + " FROM " + r.table
	if where != "" {
		q += " WHERE " + where
	}
	q += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", p.PageSize, offset)
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanRows[T](rows, r.dest, r.columns), total, nil
}

// Get returns a single entity by id.
func (r *GenericRepo[T]) Get(id int64) (*T, error) {
	q := "SELECT " + strings.Join(r.columns, ", ") + " FROM " + r.table + " WHERE id = $1"
	row := r.db.QueryRow(q, id)
	return scanRow(r.dest(), row, r.columns)
}

// Delete removes an entity by id.
func (r *GenericRepo[T]) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM "+r.table+" WHERE id = $1", id)
	return err
}

// structColumns returns the `db` tags of exported struct fields, in order.
// It always prepends "id" and appends "created_at","updated_at" when present.
func structColumns(v any) []string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	cols := []string{"id"}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		if tag == "id" || tag == "created_at" || tag == "updated_at" {
			continue
		}
		cols = append(cols, tag)
	}
	// Append audit timestamps if the struct declares them.
	has := map[string]bool{}
	for i := 0; i < t.NumField(); i++ {
		has[t.Field(i).Tag.Get("db")] = true
	}
	if has["created_at"] {
		cols = append(cols, "created_at")
	}
	if has["updated_at"] {
		cols = append(cols, "updated_at")
	}
	return cols
}

// scanRow scans a single *sql.Row into dest using the given columns.
func scanRow[T any](dest *T, row *sql.Row, cols []string) (*T, error) {
	ptrs, finalize := structPtrs(dest, cols)
	if err := row.Scan(ptrs...); err != nil {
		return nil, err
	}
	finalize()
	return dest, nil
}

// scanRows iterates rows scanning each into a fresh dest, returning the slice.
// Scan errors are logged (not silently dropped) so NULL-handling bugs surface.
func scanRows[T any](rows *sql.Rows, dest func() *T, cols []string) []T {
	out := []T{}
	for rows.Next() {
		d := dest()
		ptrs, finalize := structPtrs(d, cols)
		if err := rows.Scan(ptrs...); err != nil {
			logger.Error("row scan failed; skipping row", "columns", strings.Join(cols, ","), "error", err)
			continue
		}
		finalize()
		out = append(out, *d)
	}
	return out
}

// nullLink pairs a nullable struct pointer field with the sql.NullXxx value
// that receives its scan result, plus the kind used to dematerialize it.
type nullLink struct {
	field   reflect.Value
	scanner any // pointer to sql.NullTime / sql.NullInt64 / sql.NullFloat64 / sql.NullString
	kind    reflect.Kind
}

// structPtrs builds the scan-target pointers for the columns of struct v.
// It returns the slice of pointers (handed to rows.Scan) plus a finalize
// callback that must be invoked AFTER a successful Scan to copy nullable
// bridge values back into the struct's pointer fields. Non-pointer fields
// are scanned directly.
//
// This is what makes nullable columns (*time.Time, *int64, *float64, *string)
// safe: instead of passing a **T (which pq cannot scan NULL into), we scan
// into a sql.NullXxx and materialize the pointer after the fact.
func structPtrs(v any, cols []string) ([]any, func()) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	tagIdx := map[string]int{}
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag != "" && tag != "-" {
			tagIdx[tag] = i
		}
	}

	timeType := reflect.TypeOf(time.Time{})
	var links []nullLink
	ptrs := make([]any, 0, len(cols))

	for _, c := range cols {
		idx, ok := tagIdx[c]
		if !ok {
			var dummy sql.NullString
			ptrs = append(ptrs, &dummy)
			continue
		}
		f := rv.Field(idx)
		switch {
		case f.Kind() == reflect.Ptr && f.CanSet() && f.Type().Elem() == timeType:
			nb := &sql.NullTime{}
			ptrs = append(ptrs, nb)
			links = append(links, nullLink{field: f, scanner: nb, kind: reflect.Interface})
		case f.Kind() == reflect.Ptr && f.CanSet() && f.Type().Elem().Kind() == reflect.Int64:
			nb := &sql.NullInt64{}
			ptrs = append(ptrs, nb)
			links = append(links, nullLink{field: f, scanner: nb, kind: reflect.Int64})
		case f.Kind() == reflect.Ptr && f.CanSet() && f.Type().Elem().Kind() == reflect.Float64:
			nb := &sql.NullFloat64{}
			ptrs = append(ptrs, nb)
			links = append(links, nullLink{field: f, scanner: nb, kind: reflect.Float64})
		case f.Kind() == reflect.Ptr && f.CanSet() && f.Type().Elem().Kind() == reflect.String:
			nb := &sql.NullString{}
			ptrs = append(ptrs, nb)
			links = append(links, nullLink{field: f, scanner: nb, kind: reflect.String})
		case f.CanAddr():
			ptrs = append(ptrs, f.Addr().Interface())
		default:
			var dummy sql.NullString
			ptrs = append(ptrs, &dummy)
		}
	}

	return ptrs, func() {
		for _, ln := range links {
			switch ln.kind {
			case reflect.Interface: // *time.Time
				nb := ln.scanner.(*sql.NullTime)
				if nb.Valid {
					v := nb.Time
					ln.field.Set(reflect.ValueOf(&v))
				} else {
					ln.field.Set(reflect.Zero(ln.field.Type()))
				}
			case reflect.Int64:
				nb := ln.scanner.(*sql.NullInt64)
				if nb.Valid {
					v := nb.Int64
					ln.field.Set(reflect.ValueOf(&v))
				} else {
					ln.field.Set(reflect.Zero(ln.field.Type()))
				}
			case reflect.Float64:
				nb := ln.scanner.(*sql.NullFloat64)
				if nb.Valid {
					v := nb.Float64
					ln.field.Set(reflect.ValueOf(&v))
				} else {
					ln.field.Set(reflect.Zero(ln.field.Type()))
				}
			case reflect.String:
				nb := ln.scanner.(*sql.NullString)
				if nb.Valid {
					v := nb.String
					ln.field.Set(reflect.ValueOf(&v))
				} else {
					ln.field.Set(reflect.Zero(ln.field.Type()))
				}
			}
		}
	}
}

// Exec is a thin helper around db.Exec for custom statements.
func (r *GenericRepo[T]) Exec(query string, args ...any) (sql.Result, error) {
	return r.db.Exec(query, args...)
}
