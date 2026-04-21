# ClusterMind Architecture

You can import this directly into **Draw.io**:
1. Open [Draw.io](https://app.diagrams.net/)
2. Go to **Arrange** -> **Insert** -> **Advanced** -> **Mermaid**
3. Paste the code block below and click **Insert**

```mermaid
graph TD
    %% Define Styles
    classDef ui fill:#4a154b,stroke:#fff,stroke-width:2px,color:#fff;
    classDef core fill:#0f3460,stroke:#fff,stroke-width:2px,color:#fff;
    classDef ext fill:#16213e,stroke:#fff,stroke-width:2px,color:#fff,stroke-dasharray: 5 5;
    
    %% User Interfaces
    User([Platform Engineer / SRE])
    
    subgraph UI ["Terminal User Interface (Bubble Tea)"]
        App[Main App Loop]
        ResourcePane[Resource Pane]
        DetailPane[Detail Pane]
        ExplainPane[Explain Pane]
        RemediationPrompt[Remediation Prompt]
        QueryBar[Intent Query Bar]
        
        App --> ResourcePane
        App --> DetailPane
        App --> ExplainPane
        App --> RemediationPrompt
        App --> QueryBar
    end
    
    %% Internal Modules
    subgraph Core ["ClusterMind Internal Modules"]
        KubeClient[KubeClient Wrapper]
        ContextBuilder[Context & Evidence Builder]
        Diagnostics[Built-in Diagnostics]
        IntentEngine[Intent Parser & Executor]
        Orchestrator[AI Orchestrator]
        RemediationExec[Remediation Executor]
    end
    
    %% External Integrations
    subgraph External ["External Dependencies"]
        K8s[(Kubernetes Cluster)]
        Ollama((Local Ollama / LLM))
    end
    
    %% Relationships
    User -->|Views / Navigates| UI
    User -->|Types Query ':'| QueryBar
    User -->|Approves Fix 'Y'| RemediationPrompt
    
    %% UI to Core
    ResourcePane -->|Fetches Live Data| KubeClient
    QueryBar -->|Sends Query| IntentEngine
    ExplainPane -->|Requests RCA| ContextBuilder
    RemediationPrompt -->|Executes Action| RemediationExec
    
    %% Core Internal
    ContextBuilder -->|Pulls Events/Logs| KubeClient
    ContextBuilder -->|Runs Rules| Diagnostics
    ContextBuilder -->|Builds EvidencePackage| Orchestrator
    
    IntentEngine -->|Filters Resources| KubeClient
    RemediationExec -->|Applies Patches| KubeClient
    RemediationExec -->|Validates Actions| RemediationExec
    
    Orchestrator -->|Sends Evidence + Prompt| Ollama
    Orchestrator -->|Parses RemediationPlan| RemediationPrompt
    
    %% Core to External
    KubeClient <-->|REST API| K8s
    
    %% Apply Styles
    class UI,App,ResourcePane,DetailPane,ExplainPane,RemediationPrompt,QueryBar ui;
    class Core,KubeClient,ContextBuilder,Diagnostics,IntentEngine,Orchestrator,RemediationExec core;
    class External,K8s,Ollama ext;
```
