# Summary

We will implement an AI Chat Assistant feature that helps users discover and utilize Bucketeer's powerful but underused features. The assistant provides contextual suggestions and guidance through a conversational interface, improving feature discoverability without requiring users to read extensive documentation.

## Background

Bucketeer is a feature-rich platform offering feature flags, A/B testing, segments, progressive rollouts, flag triggers, and more. However, user research indicates that many powerful features remain underutilized:

- **Feature Discovery Problem**: Users often don't know about features like Flag Triggers, Progressive Rollout, or Prerequisite Rules until they encounter specific pain points
- **Documentation Fatigue**: Users rarely have time to read comprehensive documentation to discover relevant features
- **Context-Sensitive Needs**: Users need guidance that's relevant to their current task, not generic feature lists

This RFC proposes adding a conversational AI assistant that proactively suggests relevant features based on the user's current context within the Bucketeer dashboard.

## Goals

- Improve feature discoverability by providing contextual suggestions
- Reduce time-to-value for new and existing users
- Provide best practice recommendations based on current usage patterns
- Integrate seamlessly with existing Bucketeer architecture (Go backend, React frontend)
- Maintain security and privacy by not exposing sensitive user data to LLM
- Support both English and Japanese languages

## Non-Goals

- The AI assistant will NOT execute actions (flag changes, deployments, etc.)
- The AI assistant will NOT access or expose user attribute values or targeting data
- The AI assistant will NOT replace official documentation

## Design Overview

### User Experience

The AI Chat Assistant will be accessible via a floating help widget in the dashboard, following the Intercom/Zendesk pattern rather than a full ChatGPT-style interface:

1. **Floating Widget UI**: A small button in the bottom-right corner that opens a compact popup chat
2. **Lightweight Interaction**: Non-intrusive design that doesn't interrupt the user's workflow
3. **Context Awareness**: The assistant automatically detects the current page (Feature Flags, Targeting, Experiments, etc.)
4. **Proactive Suggestions**: Rule-based suggestions appear based on page context and flag state
5. **Conversational Interface**: Users can ask questions and receive streaming responses
6. **Preset Questions**: Quick-access buttons for common questions per page type

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│    Feature Flags Dashboard                                  │
│                                                             │
│    [my-feature-flag]                                        │
│    ├─ Status: Enabled                                       │
│    └─ Targeting: 3 rules                                    │
│                                         ┌───────────────┐   │
│                                         │ 💬 Help       │   │
│                                         │───────────────│   │
│                                         │ 💡 Tip:       │   │
│                                         │ Add Trigger   │   │
│                                         │ to this flag  │   │
│                                         │               │   │
│                                         │ [Ask question]│   │
│                                         └───────────────┘   │
│                                                    [💬]     │
└─────────────────────────────────────────────────────────────┘
```

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Frontend (React + SSE Client)                │
│  - Floating widget chat UI (Radix Popover + Tailwind)           │
│  - Page context detection                                       │
│  - Streaming response display                                   │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                    POST /v1/aichat/chat (SSE)
                    GET  /v1/aichat/suggestions
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Bucketeer API Service (Go)                      │
│  ┌──────────────────┐  ┌──────────────────┐                    │
│  │   AI Chat API    │  │   RAG Service    │                    │
│  │  - Auth/AuthZ    │  │  - Doc Retrieval │                    │
│  │  - Rate Limiting │  │  - Embeddings    │                    │
│  └────────┬─────────┘  └────────┬─────────┘                    │
│           │                     │                               │
│           └─────────┬───────────┘                               │
│                     ▼                                           │
│           ┌──────────────────┐                                 │
│           │   Chat Service   │                                 │
│           │  - Prompt Build  │                                 │
│           │  - Streaming     │                                 │
│           └────────┬─────────┘                                 │
└────────────────────┼────────────────────────────────────────────┘
                     │
                     ▼
              ┌──────────────┐
              │  OpenAI API  │
              │  gpt-4o-mini │
              └──────────────┘
```

### Key Components

| Component | Responsibility |
|-----------|----------------|
| **AI Chat API** | Authentication, authorization, rate limiting, SSE streaming |
| **Chat Service** | Prompt construction, LLM interaction, response streaming |
| **RAG Service** | Documentation retrieval using embeddings for context |
| **Suggestion Service** | Rule-based proactive suggestions (no LLM required) |

### API Design

Two primary endpoints:

1. **POST /v1/aichat/chat** (SSE Streaming)
   - Accepts conversation history and page context
   - Returns streaming text response via Server-Sent Events
   - Requires authentication (Viewer role minimum)

2. **GET /v1/aichat/suggestions**
   - Returns proactive suggestions based on current page
   - Rule-based, no LLM calls required
   - Lightweight and cacheable

### Data Privacy

| Data Type | Sent to LLM | Rationale |
|-----------|-------------|-----------|
| Flag names | Yes | Required for contextual suggestions |
| Flag descriptions | Yes | Provides additional context for better recommendations |
| Variation names | Yes | Required to understand flag structure and suggest best practices |
| Variation descriptions | Yes | Helps assistant explain flag behavior to users |
| Tag names | Yes | Required for organization recommendations |
| Rule structure | Yes | Required for complexity analysis |
| Attribute values | **No** | May contain sensitive data |
| User IDs | **No** | Personal information |
| Variation values | **No** | Business logic secrets |

### Conversation History

- **Retention**: Session-only (cleared when browser closes)
- **Storage**: Frontend state only (not persisted to backend)
- **Rationale**: Privacy-first approach, simpler compliance

## Alternatives Considered

### Alternative 1: Node.js Backend with Vercel AI SDK

**Approach**: Implement AI chat as a separate Node.js service using Vercel AI SDK

**Pros**:
- Vercel AI SDK provides excellent streaming abstractions
- Rich ecosystem of AI tooling in JavaScript
- Faster prototyping

**Cons**:
- Introduces a new runtime (Node.js) into Go-based infrastructure
- Additional operational complexity (two languages, two deployment pipelines)
- Inconsistent with existing authentication/authorization patterns
- Requires API gateway changes for routing

**Decision**: Rejected due to architectural inconsistency and operational overhead

### Alternative 2: Embedded Chat Widget (Third-Party)

**Approach**: Integrate a third-party chat solution (Intercom, Drift, etc.) with custom AI backend

**Pros**:
- Faster time to market
- Polished UI components
- Built-in analytics

**Cons**:
- Limited customization for feature flag domain knowledge
- Data leaves Bucketeer infrastructure
- Ongoing licensing costs
- Less control over AI behavior and prompts

**Decision**: Rejected due to data privacy concerns and limited domain customization

### Alternative 3: Static Documentation Bot

**Approach**: Simple retrieval-only system without LLM generation

**Pros**:
- No LLM costs
- Deterministic responses
- No hallucination risk

**Cons**:
- Cannot handle conversational follow-ups
- Poor user experience for complex questions
- Limited ability to provide contextual suggestions

**Decision**: Rejected due to poor user experience; RAG + LLM provides better balance

### Alternative 4: Full LangChain Framework

**Approach**: Use full langchaingo framework for all AI operations

**Pros**:
- Comprehensive tooling for agents, chains, memory
- Easy to add complex behaviors later

**Cons**:
- Heavy dependency for relatively simple use case
- Abstractions may hide important implementation details
- More difficult to debug and optimize

**Decision**: Partially adopted; use langchaingo selectively for RAG while using go-openai directly for chat completion

### Alternative 5: Third-Party Chat UI Libraries (assistant-ui, nlux)

**Approach**: Use a third-party React chat UI library like assistant-ui or nlux

**Pros**:
- Pre-built streaming support and auto-scrolling
- Polished UI components out of the box
- Faster initial development

**Cons**:
- ChatGPT-style full interface doesn't match "help widget" UX goal
- Additional npm dependencies to maintain
- Limited customization for Bucketeer's design system
- May require adapters for custom SSE backend

**Decision**: Rejected; implement lightweight widget using existing Radix UI + Tailwind components for better design consistency and simpler maintenance


## Security Considerations

- **Authentication**: Reuse existing `token.Verifier` for all AI Chat endpoints
- **Authorization**: Require minimum Viewer role via `checkEnvironmentRole`
- **Rate Limiting**: Per-user, per-environment, and per-organization limits to prevent abuse
- **Input Sanitization**: Maximum input length limits and basic sanitization (e.g., HTML tag removal) before sending to LLM
- **Prompt Injection Prevention**:
  - System prompt explicitly defines boundaries (e.g., refuses role changes or data extraction attempts by user input)
  - User input is wrapped with delimiters (e.g., `<user_input>` tags) to clearly separate it from system instructions
- **Error Handling**: Internal errors logged but not exposed to clients

## Cost Considerations

| Cost Factor | Mitigation |
|-------------|------------|
| OpenAI API calls | Use gpt-4o-mini (cost-effective), rate limiting |
| Embeddings | Pre-compute at build time, cache results |
| Infrastructure | Leverage existing API service, no new services |

## Future Enhancements

1. **Phase 2**: Context-aware suggestions based on flag state analysis
2. **Phase 3**: Metrics, caching, cost optimization dashboards
3. **Future**: MCP server for IDE integration (similar to DevCycle)
4. **Future**: Multi-LLM support (Anthropic, Google as alternatives)

## References

- [go-openai](https://github.com/sashabaranov/go-openai) - Go OpenAI client library
- [langchaingo](https://github.com/tmc/langchaingo) - Go LLM framework for RAG
