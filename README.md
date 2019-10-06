# Issue & Pull Request virtual assistant

Implements a [GitHub
Action](https://help.github.com/en/categories/automating-your-workflow-with-github-actions)
that performs actions on issues and/or pull requests based on configurable conditions.

At the moment it provides a single functionality to auto-label new issues and new pull requests.

## Installing

Add a `.github/workflows/main.yml` file to your repository with these
contents:

	name: Virutal Assistant

	on: [issues, pull_request]

	jobs:
	  build:

		runs-on: ubuntu-latest

		steps:
		- uses: ppapapetrou76/virtual-assistant@0.1
		  env:
			GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"

Then, add a `./github/virtual-assistant.yml` with the configuration as described
below.

## Configuration

Configuration can be stored at `./github/virtual-assistant.yml` as below

    issues:
          labels:
                - label1
                - label2
                - area:label3
          at-least-one:
                labels:
                      - priority:1
                      - priority:2
                      - priority:3
                default: priority:2
    pull-requests:
          labels:
                - label1
                - label2



For example, given this `./github/virtual-assistant.yml`:

      issues:
           labels:
                - label1
                - label2
                - area:label3
           at-least-one:
                labels:
                      - priority:1
                      - priority:2
                      - priority:3
                default: priority:2
      pull-requests:
           labels:
                - label1
                - label2

the action will 
- add to all new pull request the labels : `label1` and `label2`
- add to all new issues the labels : `label1`,`label2` and `area:label3`
- check all new issues if at least one of the labels `priority:1`,`priority:2`,`priority:3` exists and if not it will add the label `priority:2`