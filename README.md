# ollama-openai-proxy

`ollama-openai-proxy` is a proxy server that emulates the REST API of [ollama/ollama](https://github.com/ollama/ollama). Requests are forwarded to OpenAI.
The connection to OpenAI is made via a third-party [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai).
This allows you to use OpenAI models via a proxy from Jetbrains AI Assistant.

## Usage

```console
# Create .env file and fill OPENAI_API_KEY variable.
cp -p .env.local .env

docker compose up --build -d
```

### Setting Jetbrains AI Assistant

1. Open Settings.
2. Move Settings > Tools > AI Assistant > Models.
3. Check `Enable Ollama`. Check whether the `Test Connection` is successful.
4. Select the model you want to use under `Core features` and `Instant helpers`.
5. Check `Offline mode`.
6. Press `Apply`, then `OK`.
7. When you open `AI Chat`, you will be able to select a model from `Ollama`.
