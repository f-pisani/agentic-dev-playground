You are an expert Golang engineer. Your task is to implement an idiomatic Go API client based on markdown documentation 
located in the specs/ directory. 
The client should be implemented as a Go package and written to the following path: <PROVIDED_OUTPUT_PATH>.

Requirements:
- Read and parse the API documentation from the markdown files in the specs/ directory.
- Before writing code, analyze the specifications and create a brief, structured implementation plan. 
- Write this plan in a README.md file at the root of the output package.
- Implement a Go package to interact with the API described in the specs.

The client must support:
- Authentication, as described in the specs (e.g., API keys, bearer tokens).
- Pagination for endpoints that support it.
- Basic error handling: properly handle and log non-2xx HTTP responses.
- Use only the Go standard library.
- Write idiomatic Go code that is well-organized and maintainable.
- Ensure the code is valid and passes go vet without errors.
- Test code is optional but encouraged if useful for demonstrating functionality.

Instructions:
- Read the markdown specs in specs/.
- Write a clear and concise implementation plan in README.md at <PROVIDED_OUTPUT_PATH>.
- Implement the API client in Go according to that plan.
- Organize the package as you see fit (e.g., client, models, utils, etc.).
