# hawthorn

a lightweight container orchestration tool for freelancers

## about

hawthorn is a tool for running, healing & autocleaning containers that you create for your freelancing projects

cli coming soon

## roadmap

- cli
- fix issues with container reattachment

## how it works

- start the orchestrator by downloading the latest release & running it
- get your code at `http://localhost:3006/auth/login`
- verify the code
- create a new container at `http://localhost:3006/containers/new` with the schema:

  - ```json
    {
      "name": "name of the project",
      "repo_url": "link to the github repo"
    }
    ```

## benchmarks

- when running with low build volume & number of containers, hawthorn uses around 15 MB of memory and <1% CPU
