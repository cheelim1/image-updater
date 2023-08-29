# Image Updater for YAML

This action updates the imageTag key value in a YAML file in your Github repository.

This often is the case when using gitops coupled with deployment to a K8s cluster with ArgoCD & ApplicationSets.

## Example Usage
```
jobs:
  ci:
    # ... your other ci steps ...

  imageTag-updater:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Update imageTag
      uses: cheelim1/image-updater@v1.11
      with:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        REPO_NAME: 'github-repository-name'
        REPO_OWNER: 'github-repo-owner'
        FILE_PATH: 'xxx/filename.yaml'
        IMAGE_TAG: 'newImageTagHere'
```

## Inputs
### Requird Inputs
GITHUB_TOKEN = Github Token which has repository write access to the repository you want to update the file of the imageTag value.
REPO_NAME = Github Repository where the YAML file lives.
REPO_OWNER= Could be your Github Organization name or your username, depending who owns the repository.
FILE_PATH = File path of which imageTag exists to update.

### Optional Inputs
GITHUB_BRANCH= Branch name to update the imageTag in the Github repository. If not declared defaults to `main`.