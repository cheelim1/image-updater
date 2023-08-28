# Image Updater for YAML

The idea for this custom github action is to easily update the imageTag in another repository. This often is the case when using gitops coupled with deployment to a K8s cluster with ArgoCD.


### Example:
```
jobs:
  dev-ci:
    # ... your other dev-ci steps ...

  image-updater:
    needs: [dev-ci]
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Update Image Tag
      uses: cheelim1/image-updater # adjust the path to where you have stored the action
      with:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        REPO_NAME: 'github-repository-name'
        REPO_OWNER: 'owner'
        FILE_PATH: '../xxx.yaml'
        IMAGE_TAG: 'newImageTagHere'

```