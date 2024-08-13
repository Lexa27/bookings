package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("Form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("Shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some value")
	form = New(postedValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("Shows min length of 100 met when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abc123")
	form = New(postedValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("Shows min length of 1 is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("Shows an error when there should not be one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid email for non existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "example@gmail.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Form shows invalid email for a valid email address")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "bademail")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Form shows valid email when it is an invalid address")
	}

}
