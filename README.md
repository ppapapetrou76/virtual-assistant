# Issue & Pull Request virtual assistant

Implements a [GitHub
Action](https://help.github.com/en/categories/automating-your-workflow-with-github-actions)
that performs actions on issues and/or pull requests based on configurable conditions.

At the moment it provides the following actions
- Labeler
    - Auto-label issues
    - Auto-label pull requests
    - Check issues for the existence of at least one label from a given list and auto-label if it's not found
- Assigner
    - Auto-add issues to a project column - only repository projects are currently supported

## Installing

Add a `.github/workflows/main.yml` file to your repository with these
contents:

	name: Virtual Assistant

	on: [issues, pull_request]

	jobs:
	  build:

		runs-on: ubuntu-latest

		steps:
		- uses: ppapapetrou76/virtual-assistant@0.3
		  env:
			GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"

Then, add a `./github/virtual-assistant.yml` with the configuration as described
below.

## Configuration

Configuration can be stored at `./github/virtual-assistant.yml` as below

The labeler action can be configured for issues and pull-requests. 
The `labels` property accepts a list of labels and these labels will be added to the issues/pull-requests
The `actions` property accepts a list of event actions to trigger the labeler
The `at-least-one` property accepts a list of labels and a default label. 

The assigner action can be configured for issues
The `project` property is composed of a `url` property which is the url of your project (just grab it from your browser)
and a `column` property which is the name of your project column (case sensitive)


    labeler:
      issues:
        labels:
          - label1
          - label2
          - area:label3
        actions:
          - opened
          - milestoned
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
        actions:
          - opened
          - synchronize
    
    assigner:
      issues:
        project:
          url: https://github.com/ppapapetrou76/virtual-assistant/projects/1
          column: To do
        actions:
          - opened
          - milestoned




For example, given the above configuration

the action will 
- add to all new pull request the labels : `label1` and `label2`
- add to all new issues the labels : `label1`,`label2` and `area:label3`
- check all new issues if at least one of the labels `priority:1`,`priority:2`,`priority:3` exists and if not it will add the label `priority:2`
- add all new issues to the project with number `1` under the column `To do`
