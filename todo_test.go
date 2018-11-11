package todo

import (
	"testing"
)

func TestCombine(t *testing.T) {
	t.Skip("Find out, what's wrong with a symbols `[`, `]` in the project_name")
	// Prepare pattern for extraction
	re := PreparePattern("//TODO({{username}}).[{{project_name}}] {{title}}")
	if re == nil {
		t.Fatal("combined regexp must be not a nil value")
	}

	// Try to extract to-do fields from given string
	result, success, err := Extract("//TODO(john.doe).[DefaultProjectName] Fix me a drink. Please... {{custom_field}}")
	if err != nil {
		t.Fatalf("error occurred %s", err)
	}
	if !success {
		t.Fatalf("expected success")
	}
	if result == nil {
		t.Fatal("result should be not (nil")
	}

	// Expected results
	expectedUsername := "john.doe"
	expectedProjectName := "{DefaultProjectName}"
	expectedIssueTitle := "Fix me a drink. Please... {{custom_field}}"

	if result.Username != expectedUsername {
		t.Errorf("result.Username `%s` != `%s`", expectedUsername, result.Username)
	}
	if result.ProjectName != expectedProjectName {
		t.Errorf("result.ProjectName `%s` != `%s`", expectedProjectName, result.ProjectName)
	}
	if result.IssueTitle != expectedIssueTitle {
		t.Errorf("result.IssueTitle `%s` != `%s`", expectedIssueTitle, result.IssueTitle)
	}
}
