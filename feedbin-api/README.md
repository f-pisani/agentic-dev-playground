# Feedbin API Client Implementation: A Developer's Comparison of AI Tools

*Note: This document was AI-generated based on developer feedback and prompting, summarizing personal experiences with various AI coding tools.*

## Introduction
Feedbin is a web-based RSS reader that offers both a user interface for managing feeds and a REST-like API for client applications. As a software developer exploring ways to enhance my workflow, I decided to test various AI coding assistants by having them implement a Go client for the Feedbin API. This document captures my personal observations and experiences rather than a formal benchmark or research study.

## Experiment Setup
I tasked multiple AI coding assistants with implementing a Go client library for the Feedbin API. Each assistant received identical instructions and had access only to the official API documentation (from the README at https://github.com/feedbin/feedbin-api). The assistants worked independently and could not reference each other's implementations.

My primary interest was in:
- Whether the assistant could complete the implementation in one continuous session
- The presence of errors in the generated code
- Whether the implementation worked when tested
- If the assistant provided a working example
- The overall developer experience

## Implementation Completeness

After further inspection, it's important to note that none of the AI tools managed to implement ALL Feedbin API endpoints in one go. While many implemented a substantial portion of the API, certain endpoints were consistently missed across implementations, including:

- Icons endpoints
- Imports endpoints
- Pages endpoints
- Saved Searches endpoints

This limitation affected even the tools marked as "Complete success" below. The success ratings are relative to each other rather than measuring against full API coverage.

## Results by AI Tool

### Claude-Based Tools

#### Claude Code CLI (~$1)
**Outcome**: Mostly complete with working code
- Implemented many API endpoints in one continuous session
- Code had no errors and worked properly for implemented endpoints
- Included a working example
- Missing some endpoints (icons, imports, pages, saved searches)
- Cost was approximately $1 for API usage

#### JetBrains Junie (Claude 3.7)
**Outcome**: Mostly complete with working code
- Delivered a functional implementation in one continuous session
- No errors in the generated code
- Included a working example
- Missing some endpoints (icons, imports, pages, saved searches)

#### Cursor Claude 3.7
**Outcome**: Mostly complete with working code
- Completed many API endpoints in one continuous session
- Generated error-free code
- Provided a working example
- Missing some endpoints (icons, imports, pages, saved searches)

#### VS Code RooCode in Architect mode (Claude 3.7) (~$1.6)
**Outcome**: Mostly complete with working code
- Implemented many API endpoints in one continuous session
- Code had no errors and worked correctly for implemented endpoints
- Included a working example
- Missing some endpoints (icons, imports, pages, saved searches)
- Cost was approximately $1.6 for API usage (with Architect mode)

#### Aider Claude 3.7
**Outcome**: Partial success with issues
- Could not complete the implementation in one session
- The Go module wasn't initialized correctly
- Generated code contained errors
- The provided example didn't work as expected
- Missing multiple endpoints, including the commonly missed ones

#### JetBrains Augment Claude 3.7
**Outcome**: Partial success with interruptions
- Implementation was interrupted midway, requiring manual prompting to continue
- Missing several endpoints, including the commonly missed ones
- Final result included working code and example for implemented endpoints

#### JetBrains Cascade Claude 3.7
**Outcome**: Partial success with interruptions
- Similar to Augment, the implementation halted mid-generation
- Required prompting to continue
- Missing several endpoints, including the commonly missed ones
- Final code worked correctly with a functional example for implemented endpoints

#### Windsurf (Claude 3.7)
**Outcome**: Partial success with interruptions
- Implementation process was interrupted
- Required manual intervention to resume generation
- Missing several endpoints, including the commonly missed ones
- Final result included working code and example for implemented endpoints

#### Trae (Claude 3.7)
**Outcome**: Failure
- Failed to deliver a complete implementation
- Repeatedly explained what needed to be done without generating the actual code
- Did not produce a functional solution
- Many endpoints not implemented at all

### Other Models

#### Cursor Gemini 2.5 Pro Exp 03 25
**Outcome**: Incomplete implementation
- Stopped generating code partway through the task
- Partial implementation had errors and missing endpoints
- Did not provide an example
- Overall implementation was non-functional
- Many endpoints not implemented, including the commonly missed ones

#### JetBrains Cascade Gemini 2.5 Pro Exp 03 25
**Outcome**: Incomplete implementation with interruptions
- Halted mid-implementation, requiring manual prompting
- Even after continuation, no working example was provided
- Implementation was incomplete with many missing endpoints
- Commonly missed endpoints were not implemented

#### JetBrains Cascade GPT 4.1
**Outcome**: Incomplete implementation with interruptions
- Stopped generating code before completion
- Required manual prompting to continue
- No working example was provided
- Many endpoints not implemented, including the commonly missed ones

#### JetBrains Cascade Deepseek R1
**Outcome**: Complete failure
- Failed to generate relevant code
- Produced Python code instead of Go
- Occasionally responded in Chinese
- Could not deliver any usable implementation
- No endpoints properly implemented

## Key Observations

As a software developer exploring these tools to enhance my workflow, I noticed several patterns:

1. **Endpoint Coverage Limitations**: No tool implemented the complete Feedbin API in one go. Even the best performers consistently missed certain endpoints (icons, imports, pages, saved searches).

2. **Performance by Model**: Tools based on Claude 3.7 generally delivered more complete and functional implementations compared to others. The most reliable performers were Claude Code CLI, JetBrains Junie, Cursor Claude 3.7, and VS Code RooCode.

3. **Completion Reliability**: Many tools had difficulty maintaining focus through the entire implementation. Some stopped mid-generation and required prompting to continue, while others (like Trae) never produced complete solutions despite being based on capable models.

4. **Practical Considerations**: 
   - Cost efficiency varied, with some effective implementations available at reasonable API usage costs ($1-1.6)
   - Tools with better integration into development environments provided a smoother experience when handling interruptions

5. **Integration Matters**: Even with the same underlying model (Claude 3.7), different tools produced significantly different results, suggesting that how a model is integrated into a development environment strongly affects its performance.

## Personal Takeaways

As a developer looking to incorporate AI assistance into my workflow, I found that:

- Claude 3.7-based tools currently offer the best balance of quality and reliability for implementing API clients in Go, though complete coverage remains a challenge
- For API client implementations, several tools can deliver working code in a single session, but manual review and additions are still necessary
- The most effective tools for my workflow were JetBrains Junie, Cursor Claude 3.7, Claude Code CLI, and VS Code RooCode
- The specific integration of an AI model into a development environment significantly impacts its usefulness
- AI tools are valuable for quickly scaffolding API clients, but developers should expect to manually complete missing endpoints

This informal comparison has helped me understand which AI coding assistants might best complement my development process for tasks similar to implementing API clients, while also clarifying their current limitations.