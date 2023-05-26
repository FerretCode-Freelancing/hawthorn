# hawthorn

a lightweight container orchestration tool for freelancers

## about

hawthorn is a tool for running, healing & autocleaning containers that you create for your freelancing projects

## how it works

- start the orchestrator by downloading the latest release & running it
- authenticate at `http://localhost:3006/auth/login`
- create a new container at `http://localhost:3006/containers/new` with the schema:

  - ```json
    {
      "name": "name of the project",
      "repo_url": "link to the github repo"
    }
    ```
