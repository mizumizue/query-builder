package object_parser

import (
	"reflect"
	"testing"
	"time"
)

type Date struct {
	time.Time
}

type SearchArticleRequest struct {
	Title         *string `query:"title" search:"title" operator:"eq"`
	PublishedDate *Date   `query:"publishedDate" search:"published_date" operator:"ge"`
}

func allocateStr(value string) *string {
	return &value
}

func allocateDate(value Date) *Date {
	return &value
}

var now = time.Now()
var sar1 = SearchArticleRequest{
	Title:         allocateStr("TestTitle"),
	PublishedDate: allocateDate(Date{Time: now}),
}
var sar2 = SearchArticleRequest{
	Title: allocateStr("TestTitle"),
}
var sar3 = SearchArticleRequest{
	PublishedDate: allocateDate(Date{Time: now}),
}

func TestParameterParser_ParseNamedParam(t *testing.T) {
	type fields struct {
		object      interface{}
		objectType  reflect.Type
		objectValue reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "",
			fields: fields{
				object:      sar1,
				objectType:  reflect.TypeOf(sar1),
				objectValue: reflect.ValueOf(sar1),
			},
			want: map[string]interface{}{
				"title":          "TestTitle",
				"published_date": Date{now},
			},
		},
		{
			name: "",
			fields: fields{
				object:      sar2,
				objectType:  reflect.TypeOf(sar2),
				objectValue: reflect.ValueOf(sar2),
			},
			want: map[string]interface{}{
				"title": "TestTitle",
			},
		},
		{
			name: "",
			fields: fields{
				object:      sar3,
				objectType:  reflect.TypeOf(sar3),
				objectValue: reflect.ValueOf(sar3),
			},
			want: map[string]interface{}{
				"published_date": Date{now},
			},
		},
		{
			name: "",
			fields: fields{
				object:      SearchArticleRequest{},
				objectType:  reflect.TypeOf(SearchArticleRequest{}),
				objectValue: reflect.ValueOf(SearchArticleRequest{}),
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectParser := &ObjectParser{
				object:      tt.fields.object,
				objectType:  tt.fields.objectType,
				objectValue: tt.fields.objectValue,
			}
			if got := objectParser.NamedParam(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectParser.NamedParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParameterParser_ParseSearchBindMap(t *testing.T) {
	type fields struct {
		object      interface{}
		objectType  reflect.Type
		objectValue reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]string
	}{
		{
			name: "",
			fields: fields{
				object:      sar1,
				objectType:  reflect.TypeOf(sar1),
				objectValue: reflect.ValueOf(sar1),
			},
			want: []map[string]string{
				{
					"bind":     "title",
					"target":   "title",
					"operator": "eq",
				},
				{
					"bind":     "published_date",
					"target":   "published_date",
					"operator": "ge",
				},
			},
		},
		{
			name: "",
			fields: fields{
				object:      SearchArticleRequest{},
				objectType:  reflect.TypeOf(SearchArticleRequest{}),
				objectValue: reflect.ValueOf(SearchArticleRequest{}),
			},
			want: []map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectParser := &ObjectParser{
				object:      tt.fields.object,
				objectType:  tt.fields.objectType,
				objectValue: tt.fields.objectValue,
			}
			if got := objectParser.SearchBindMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectParser.SearchBindMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
