/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootGet(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	rootGet(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthzGet(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	healthzGet(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerFunc(t *testing.T) {
	t.Parallel()

}
