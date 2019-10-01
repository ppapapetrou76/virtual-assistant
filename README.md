# Issue & Pull Request virtual assistant

Implements a [GitHub
Action](https://help.github.com/en/categories/automating-your-workflow-with-github-actions)
that performs actions on Issues and/or Pull Requests based on configurable conditions.

At the moment it provides a single functionality to auto-label all new pull requests


## Installing

Add a .github/workflows/main.yml file to your repository with these
contents:

	name: Virutal Assistant

	on:
	  - pull_request

	jobs:
	  build:

		runs-on: ubuntu-latest
		
		steps:
		- uses: ppapapetrou76/virtual-assistant@master
		  env:
			GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"

Then, add a `./github/virtual-assistant.yml` with the configuration as described
below.

## Configuration

Configuration can be stored at `./github/virtual-assistant.yml` as a plain list of labels
    
    labels:
      - <label1>
      - <label2>
      

For example, given this `./github/virtual-assistant.yml`:

      labels:
            - label1
            - label2
            - area:label3

