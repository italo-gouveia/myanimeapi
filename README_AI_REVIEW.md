# 🤖 Revisão Automática de Código com IA

Este projeto agora possui revisão automática de código usando Inteligência Artificial (ChatGPT) através do GitHub Actions. A IA analisa automaticamente todos os pull requests e fornece comentários detalhados sobre qualidade, segurança, performance e boas práticas.

## ✨ Funcionalidades

- 🔍 **Análise Automática**: Revisão automática em todos os PRs
- 🔒 **Foco em Segurança**: Detecção de vulnerabilidades e problemas de segurança
- ⚡ **Análise de Performance**: Sugestões de otimização
- 📝 **Qualidade de Código**: Verificação de boas práticas de Go
- 🏗️ **Arquitetura**: Análise de design patterns e estrutura
- 🌐 **APIs**: Validação específica para APIs REST

## 🚀 Configuração Rápida

### 1. Obter Chave da OpenAI

1. Acesse [OpenAI Platform](https://platform.openai.com/api-keys)
2. Crie uma nova API Key
3. Copie a chave (começa com `sk-`)

### 2. Configurar no GitHub

1. Vá para seu repositório → Settings → Secrets and variables → Actions
2. Clique em "New repository secret"
3. Nome: `OPENAI_API_KEY`
4. Valor: Cole sua chave da OpenAI

### 3. Testar

1. Crie um Pull Request
2. Aguarde 2-5 minutos
3. Verifique os comentários da IA no PR

## 📁 Arquivos Criados

```
.github/
├── workflows/
│   ├── code-review.yml          # Workflow principal de revisão
│   └── security-review.yml      # Workflow específico de segurança
├── pr-agent.toml               # Configurações do PR Agent
└── .golangci.yml               # Configuração do linter

docs/
├── AI_CODE_REVIEW.md           # Documentação completa
└── GITHUB_SECRETS_SETUP.md     # Guia de configuração de secrets
```

## 💰 Custos

- **GPT-3.5-turbo**: ~$0.002 por 1K tokens
- **GPT-4**: ~$0.03 por 1K tokens
- **Revisão típica**: $0.01-$0.10 por PR

## 🎯 Exemplo de Comentário da IA

```
🔧 **Sugestão de Melhoria**

Na função `CreateAnime`, considere adicionar validação de entrada:

```go
// Antes
func (s *AnimeService) CreateAnime(anime *models.Anime) error {
    return s.repo.Create(anime)
}

// Depois
func (s *AnimeService) CreateAnime(anime *models.Anime) error {
    if anime == nil {
        return errors.New("anime cannot be nil")
    }
    
    if anime.Title == "" {
        return errors.New("anime title is required")
    }
    
    return s.repo.Create(anime)
}
```

🔒 **Segurança**: Sempre valide dados de entrada para prevenir injeção de dados maliciosos.
```

## 🔧 Configurações Avançadas

### Personalizar Instruções

Edite `.github/pr-agent.toml`:

```toml
[pr_reviewer]
extra_instructions = """
Suas instruções personalizadas aqui...
"""
```

### Ignorar Arquivos

```toml
ignore_files = [
    "*.md",
    "docs/**",
    "scripts/**"
]
```

### Mudar Modelo

```toml
model = "gpt-3.5-turbo"  # Mais rápido, menos preciso
model = "gpt-4"          # Mais preciso, mais lento
```

## 🛠️ Workflows Disponíveis

### 1. AI Code Review (Principal)
- **Arquivo**: `.github/workflows/code-review.yml`
- **Trigger**: PRs para branch `main`
- **Função**: Revisão geral de código

### 2. Security Review (Específico)
- **Arquivo**: `.github/workflows/security-review.yml`
- **Trigger**: PRs para branch `main`
- **Função**: Análise focada em segurança

## 📊 Categorias de Análise

| Emoji | Categoria | Descrição |
|-------|-----------|-----------|
| 🔧 | Melhoria | Sugestões de melhorias gerais |
| 🐛 | Bug | Problemas identificados |
| ⚠️ | Aviso | Pontos de atenção |
| 💡 | Sugestão | Ideias e recomendações |
| 🔒 | Segurança | Vulnerabilidades e riscos |
| ⚡ | Performance | Otimizações possíveis |
| 📝 | Documentação | Melhorias na documentação |
| 🧪 | Testes | Sugestões de testes |

## 🚨 Troubleshooting

### Problema: Revisão não funciona
**Solução**: Verifique se `OPENAI_API_KEY` está configurado corretamente

### Problema: Comentários muito genéricos
**Solução**: Ajuste as instruções em `.github/pr-agent.toml`

### Problema: Muitos falsos positivos
**Solução**: Diminua a temperatura ou refine as instruções

### Problema: Workflow não executa
**Solução**: Verifique se os arquivos estão na branch correta

## 📚 Documentação Completa

- [Guia Detalhado](docs/AI_CODE_REVIEW.md) - Documentação completa
- [Configuração de Secrets](docs/GITHUB_SECRETS_SETUP.md) - Como configurar secrets

## 🤝 Contribuindo

Para melhorar a revisão automática:

1. **Reporte problemas** em issues
2. **Sugira melhorias** nas configurações
3. **Teste diferentes parâmetros**
4. **Atualize a documentação**

## ⚖️ Limitações

- **Não substitui revisão humana**: Use como complemento
- **Pode ter falsos positivos**: Sempre use seu julgamento
- **Depende da qualidade das instruções**: Configure adequadamente
- **Custo da API**: Cada revisão consome tokens

## 🎉 Benefícios

- ✅ **Feedback imediato** em todos os PRs
- ✅ **Consistência** nos padrões de revisão
- ✅ **Aprendizado** contínuo da equipe
- ✅ **Qualidade** de código melhorada
- ✅ **Produtividade** aumentada

---

**Pronto para usar!** 🚀

Após configurar o secret `OPENAI_API_KEY`, a revisão automática estará ativa em todos os seus pull requests.
