package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

func extractBranchName(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return "main" // default to "main" if extraction fails
}

func main() {
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	if token == "" {
		log.Fatal("INPUT_GITHUB_TOKEN is not set")
	}

	repoOwner := os.Getenv("INPUT_REPO_OWNER")
	if repoOwner == "" {
		log.Fatal("INPUT_REPO_OWNER is not set")
	}

	repoName := os.Getenv("INPUT_REPO_NAME")
	if repoName == "" {
		log.Fatal("INPUT_REPO_NAME is not set")
	}

	filePath := os.Getenv("INPUT_FILE_PATH")
	if filePath == "" {
		log.Fatal("INPUT_FILE_PATH is not set")
	}

	newTag := os.Getenv("INPUT_IMAGE_TAG")
	if newTag == "" {
		log.Fatal("INPUT_IMAGE_TAG is not set")
	}

	branch := extractBranchName(os.Getenv("INPUT_GITHUB_BRANCH"))
	if branch == "" {
		branch = "main" // default to "main" if INPUT_GITHUB_BRANCH is not set
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, filePath, nil)
	if err != nil {
		log.Fatalf("Failed to get content: %v", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		log.Fatalf("Failed to decode content: %v", err)
	}

	var yamlData map[interface{}]interface{}
	err = yaml.Unmarshal([]byte(content), &yamlData)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Deep-search and update "imageTag" key.
	updateImageTag(yamlData, "imageTag", newTag)

	updatedYaml, err := yaml.Marshal(yamlData)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	encodedContent := base64.StdEncoding.EncodeToString(updatedYaml)
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(fmt.Sprintf("Update imageTag to %s", newTag)),
		Content:   []byte(encodedContent),
		SHA:       fileContent.SHA,
		Branch:    github.String(branch),
		Committer: &github.CommitAuthor{Name: github.String("GitHub Actions"), Email: github.String("actions@github.com")},
	}

	_, _, err = client.Repositories.UpdateFile(ctx, repoOwner, repoName, filePath, opts)
	if err != nil {
		log.Fatalf("Failed to update file: %v", err)
	}

	fmt.Println("Image tag updated successfully!")
}

// Recursively searches and updates the given key in a nested map.
func updateImageTag(data map[interface{}]interface{}, key string, newValue string) {
	if v, ok := data[key]; ok && v != nil {
		data[key] = newValue
	}
	for _, v := range data {
		if nestedMap, ok := v.(map[interface{}]interface{}); ok {
			updateImageTag(nestedMap, key, newValue)
		}
	}
}
