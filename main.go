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
    // Fetch the necessary environment variables based on action.yml
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

    // Setup the GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

    // Fetch the file content from the GitHub repo
	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, filePath, nil)
	if err != nil {
		log.Fatalf("Failed to get content: %v", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		log.Fatalf("Failed to decode content: %v", err)
	}

    // Update the 'imageTag' in the YAML
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

    // Update the file in the GitHub repo with the new content
	encodedContent := base64.StdEncoding.EncodeToString(updatedYaml)
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(fmt.Sprintf("Update imageTag to %s", newTag)),
		Content:   []byte(encodedContent),
		SHA:       fileContent.SHA,
		Branch:    github.String("main"), // or whatever branch you want
		Committer: &github.CommitAuthor{Name: github.String("GitHub Actions"), Email: github.String("actions@github.com")},
	}

	_, _, err = client.Repositories.UpdateFile(ctx, repoOwner, repoName, filePath, opts)
	if err != nil {
		log.Fatalf("Failed to update file: %v", err)
	}

	fmt.Println("Image tag updated successfully!")
}
