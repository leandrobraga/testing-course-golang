package main

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it should not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_required(t *testing.T) {
	req := httptest.NewRequest("POST", "http://testing", nil)
	form := NewForm(req.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("forms shows valid when required fileds are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r := httptest.NewRequest("POST", "http://testing1", nil)
	r.PostForm = postedData

	form = NewForm(r.PostForm)

	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Errorf("shows post does not have required fields, when it does")
	}
}

func TestForm_Check(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("Valid() returns false, and it should be true when calling Check()")
	}
}

func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("should have an error returned from Get, but do not")
	}

	s = form.Errors.Get("whatever")

	if len(s) != 0 {
		t.Error("should not have an error returned from Get, but got one")
	}
}