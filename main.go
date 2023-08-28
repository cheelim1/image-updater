package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v37/github"
	"gopkg.in/yaml.v2"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	repoName := os.Getenv("REPO_NAME")
	if repoName == "" {
		log.Fatal("REPO_NAME is not set")
	}

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		log.Fatal("FILE_PATH is not set")
	}

	newTag := os.Getenv("IMAGE_TAG")
	if newTag == "" {
		log.Fatal("IMAGE_TAG is not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, "ml", repoName, filePath, nil)
	if err != nil {
		log.Fatalf("Failed to get content: %v", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		log.Fatalf("Failed to decode content: %v", err)
	}

	var yamlData map[string]interface{}
	err = yaml.Unmarshal([]byte(content), &yamlData)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	yamlData["imageTag"] = newTag
	updatedYaml, err := yaml.Marshal(yamlData)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	encodedContent := base64.StdEncoding.EncodeToString(updatedYaml)
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(fmt.Sprintf("Update imageTag to %s", newTag)),
		Content:   []byte(encodedContent),
		SHA:       fileContent.SHA,
		Branch:    github.String("main"), // or whatever branch you want
		Committer: &github.CommitAuthor{Name: github.String("GitHub Actions"), Email: github.String("actions@github.com")},
	}

	_, _, err = client.Repositories.UpdateFile(ctx, "ml", repoName, filePath, opts)
	if err != nil {
		log.Fatalf("Failed to update file: %v", err)
	}

	fmt.Println("Image tag updated successfully!")
}

