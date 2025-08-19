# Configuração de Secrets no GitHub

Este documento explica como configurar os secrets necessários para que a revisão automática de código funcione corretamente.

## Secrets Necessários

### 1. OPENAI_API_KEY

Este é o secret mais importante para a revisão automática de código funcionar.

#### Como obter:

1. **Acesse a OpenAI:**
   - Vá para [OpenAI Platform](https://platform.openai.com/)
   - Faça login ou crie uma conta

2. **Crie uma API Key:**
   - Clique em "API Keys" no menu lateral
   - Clique em "Create new secret key"
   - Dê um nome descritivo (ex: "GitHub Code Review")
   - Copie a chave gerada (começa com `sk-`)

3. **Adicione no GitHub:**
   - Vá para seu repositório no GitHub
   - Clique em "Settings" (aba superior)
   - No menu lateral, clique em "Secrets and variables" > "Actions"
   - Clique em "New repository secret"
   - Nome: `OPENAI_API_KEY`
   - Valor: Cole a chave da OpenAI (ex: `sk-abc123...`)

#### Custo:

- **GPT-3.5-turbo**: ~$0.002 por 1K tokens
- **GPT-4**: ~$0.03 por 1K tokens
- Uma revisão típica custa entre $0.01-$0.10

### 2. GITHUB_TOKEN (Opcional)

Este token é automaticamente fornecido pelo GitHub Actions, mas você pode configurar um personalizado se necessário.

#### Como configurar (se necessário):

1. **Crie um Personal Access Token:**
   - Vá para GitHub Settings > Developer settings > Personal access tokens
   - Clique em "Tokens (classic)"
   - Clique em "Generate new token"
   - Selecione os escopos: `repo`, `pull_requests`
   - Copie o token gerado

2. **Adicione no GitHub:**
   - Nome: `GITHUB_TOKEN`
   - Valor: Cole o token personalizado

## Configuração no GitHub

### Passo a Passo:

1. **Acesse o repositório:**
   ```
   https://github.com/seu-usuario/myanimeapi
   ```

2. **Vá para Settings:**
   - Clique na aba "Settings" (ao lado de "Code", "Issues", etc.)

3. **Acesse Secrets:**
   - No menu lateral, clique em "Secrets and variables"
   - Clique em "Actions"

4. **Adicione o secret:**
   - Clique em "New repository secret"
   - Nome: `OPENAI_API_KEY`
   - Valor: Sua chave da OpenAI
   - Clique em "Add secret"

### Verificação:

Para verificar se está funcionando:

1. **Crie um Pull Request** de teste
2. **Aguarde 2-5 minutos** para a análise
3. **Verifique os comentários** no PR
4. **Verifique os logs** em Actions > AI Code Review

## Troubleshooting

### "Secret not found" error:

- Verifique se o nome do secret está exatamente como `OPENAI_API_KEY`
- Certifique-se de que o secret foi adicionado no repositório correto
- Verifique se você tem permissão para acessar secrets

### "Invalid API key" error:

- Verifique se a chave da OpenAI está correta
- Certifique-se de que a chave não expirou
- Verifique se você tem créditos na conta da OpenAI

### Workflow não executa:

- Verifique se os arquivos `.github/workflows/code-review.yml` existem
- Verifique se o workflow está configurado para a branch correta
- Verifique se você tem permissão para executar workflows

### Comentários não aparecem:

- Verifique se o workflow tem permissão para comentar em PRs
- Verifique os logs do workflow para erros
- Certifique-se de que o PR está na branch `main`

## Segurança

### Boas Práticas:

1. **Nunca commite secrets** no código
2. **Use secrets do repositório** em vez de secrets da organização quando possível
3. **Rotacione as chaves** periodicamente
4. **Monitore o uso** da API da OpenAI
5. **Configure limites de gastos** na OpenAI

### Configuração de Limites:

Na OpenAI Platform, você pode configurar:

- **Limite de gastos** por mês
- **Limite de requisições** por minuto
- **Alertas** quando atingir limites

## Exemplo de Configuração Completa

```yaml
# .github/workflows/code-review.yml
name: AI Code Review

on:
  pull_request:
    types: [opened, synchronize, reopened]
    branches:
      - main

permissions:
  contents: read
  pull-requests: write
  actions: read

jobs:
  ai-code-review:
    name: AI Code Review
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Run AI Code Review
        uses: Codium-ai/pr-agent@main
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          openai_api_key: ${{ secrets.OPENAI_API_KEY }}
          mode: review
          review_comment_lgtm: false
          review_suggestions: true
          review_estimate: true
```

## Próximos Passos

Após configurar os secrets:

1. **Teste o sistema** criando um PR de teste
2. **Ajuste as configurações** conforme necessário
3. **Monitore os custos** da API
4. **Treine sua equipe** sobre como usar os comentários da IA
5. **Refine as instruções** baseado no feedback recebido
