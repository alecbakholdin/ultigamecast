package validation

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v5"
)

func AddFieldErrorString(c echo.Context, key string, s string) {
	AddFieldError(c, key, fmt.Errorf(s))
}

func AddFieldError(c echo.Context, key string, err error) {
	c.Set("invalid", true)
	m := getErrorMap(c)
	if e, ok := m[key]; ok {
		m[key] = append(e, err.Error())
	} else {
		m[key] = []string{err.Error()}
	}
}

func IsFieldValid(c echo.Context, key string) bool {
	m := getErrorMap(c)
	_, ok := m[key]
	return !ok
}

func GetFieldErrorString(c echo.Context, key string) string {
	m := getErrorMap(c)
	if e, ok := m[key]; ok {
		return strings.Join(e, ", ")
	}
	return ""
}

func getErrorMap(c echo.Context) map[string][]string {
	if m, ok := c.Get("errors").(map[string][]string); ok {
		return m
	} else {
		m = make(map[string][]string)
		c.Set("errors", m)
		return m
	}
}

func AddFormErrorString(c echo.Context, s string) {
	AddFormError(c, fmt.Errorf(s))
}

func AddFormError(c echo.Context, err error) {
	c.Set("invalid", true)
	if f, ok := c.Get("form_error").([]string); ok {
		c.Set("form_error", []string{err.Error()})
	} else {
		c.Set("form_error", append(f, err.Error()))
	}
}

func GetFormErrorString(c echo.Context) string {
	if f, ok := c.Get("form_error").([]string); ok {
		return strings.Join(f, ", ")
	} else {
		return ""
	}
}

func IsFormValid(c echo.Context) bool {
	invalid, ok := c.Get("invalid").(bool)
	if !ok {
		return true
	} else {
		return !invalid
	}
}
