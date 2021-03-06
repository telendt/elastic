// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestIndicesExistsTypeBuildURL(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tests := []struct {
		Indices               []string
		Types                 []string
		Expected              string
		ExpectValidateFailure bool
	}{
		{
			[]string{},
			[]string{},
			"",
			true,
		},
		{
			[]string{"index1"},
			[]string{},
			"",
			true,
		},
		{
			[]string{},
			[]string{"type1"},
			"",
			true,
		},
		{
			[]string{"index1"},
			[]string{"type1"},
			"/index1/_mapping/type1",
			false,
		},
		{
			[]string{"index1", "index2"},
			[]string{"type1"},
			"/index1%2Cindex2/_mapping/type1",
			false,
		},
		{
			[]string{"index1", "index2"},
			[]string{"type1", "type2"},
			"/index1%2Cindex2/_mapping/type1%2Ctype2",
			false,
		},
	}

	for i, test := range tests {
		err := client.TypeExists().Index(test.Indices...).Type(test.Types...).Validate()
		if err == nil && test.ExpectValidateFailure {
			t.Errorf("#%d: expected validate to fail", i+1)
			continue
		}
		if err != nil && !test.ExpectValidateFailure {
			t.Errorf("#%d: expected validate to succeed", i+1)
			continue
		}
		if !test.ExpectValidateFailure {
			path, _, err := client.TypeExists().Index(test.Indices...).Type(test.Types...).buildURL()
			if err != nil {
				t.Fatalf("#%d: %v", i+1, err)
			}
			if path != test.Expected {
				t.Errorf("#%d: expected %q; got: %q", i+1, test.Expected, path)
			}
		}
	}
}

func TestIndicesExistsType(t *testing.T) {
	client := setupTestClient(t)

	// Create index with tweet type
	createIndex, err := client.CreateIndex(testIndexName).Body(testMapping).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}
	if !createIndex.Acknowledged {
		t.Errorf("expected CreateIndexResult.Acknowledged %v; got %v", true, createIndex.Acknowledged)
	}

	// Check if type exists
	exists, err := client.TypeExists().Index(testIndexName).Type("_doc").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("type %s should exist in index %s, but doesn't\n", "_doc", testIndexName)
	}

	// Delete index
	deleteIndex, err := client.DeleteIndex(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !deleteIndex.Acknowledged {
		t.Errorf("expected DeleteIndexResult.Acknowledged %v; got %v", true, deleteIndex.Acknowledged)
	}

	// Check if type exists
	exists, err = client.TypeExists().Index(testIndexName).Type("doc").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("type %s should not exist in index %s, but it does\n", "_doc", testIndexName)
	}
}

func TestIndicesExistsTypeValidate(t *testing.T) {
	client := setupTestClient(t)

	// No index name -> fail with error
	res, err := NewIndicesExistsTypeService(client).Do(context.TODO())
	if err == nil {
		t.Fatalf("expected IndicesExistsType to fail without index name")
	}
	if res != false {
		t.Fatalf("expected result to be false; got: %v", res)
	}
}
