# Revisão Automática de Código com IA

Este projeto utiliza revisão automática de código com Inteligência Artificial através do GitHub Actions para analisar pull requests e fornecer feedback construtivo.

## Como Funciona

A revisão automática é executada sempre que um pull request é:
- Aberto
- Atualizado (synchronize)
- Reaberto

O sistema utiliza o **PR Agent** da CodiumAI para analisar as mudanças e fornecer comentários detalhados sobre:

- 🔧 **Qualidade do Código**: Boas práticas de Go, estrutura, nomenclatura
- 🔒 **Segurança**: Vulnerabilidades, validação de entrada, controle de acesso
- ⚡ **Performance**: Otimizações, uso de recursos, padrões de concorrência
- 📝 **Manutenibilidade**: Legibilidade, documentação, testes
- 🏗️ **Arquitetura**: Design patterns, separação de responsabilidades
- 🌐 **APIs**: Validação, tratamento de erros HTTP, middleware

## Configuração

### 1. Configurar Secrets no GitHub

No seu repositório GitHub, vá para **Settings > Secrets and variables > Actions** e adicione:

```
OPENAI_API_KEY=sk-your-openai-api-key-here
```

### 2. Obter Chave da OpenAI

1. Acesse [OpenAI API Keys](https://platform.openai.com/api-keys)
2. Crie uma nova chave de API
3. Copie a chave e adicione como secret no GitHub

### 3. Configurações do Projeto

O arquivo `.github/pr-agent.toml` contém as configurações específicas para o projeto:

- **Modelo**: GPT-4 para análises mais precisas
- **Tokens**: 4000 para análises detalhadas
- **Arquivos ignorados**: Documentação, dependências, frontend
- **Instruções específicas**: Focadas em Go e APIs

## Como Usar

### Para Desenvolvedores

1. **Crie um Pull Request** normalmente
2. **Aguarde a análise automática** (geralmente 2-5 minutos)
3. **Revise os comentários** da IA no PR
4. **Implemente as sugestões** quando apropriado
5. **Responda aos comentários** se necessário

### Exemplo de Comentário da IA

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

## Configurações Avançadas

### Personalizar Instruções

Edite o arquivo `.github/pr-agent.toml` para personalizar:

```toml
[pr_reviewer]
extra_instructions = """
Suas instruções personalizadas aqui...
"""
```

### Ignorar Arquivos

Adicione padrões ao array `ignore_files`:

```toml
ignore_files = [
    "*.md",
    "docs/**",
    "scripts/**"
]
```

### Configurar Modelo

Altere o modelo usado para análise:

```toml
model = "gpt-3.5-turbo"  # Mais rápido, menos preciso
model = "gpt-4"          # Mais preciso, mais lento
```

## Troubleshooting

### A revisão não está funcionando

1. **Verifique os secrets**: Certifique-se de que `OPENAI_API_KEY` está configurado
2. **Verifique as permissões**: O workflow precisa de permissão para comentar em PRs
3. **Verifique os logs**: Acesse a aba Actions no GitHub para ver logs detalhados

### Comentários muito genéricos

1. **Ajuste as instruções**: Modifique `extra_instructions` no `.github/pr-agent.toml`
2. **Use um modelo melhor**: Mude para `gpt-4` se estiver usando `gpt-3.5-turbo`
3. **Aumente os tokens**: Incremente `max_tokens` para análises mais detalhadas

### Muitos falsos positivos

1. **Ajuste a temperatura**: Diminua `temperature` para respostas mais conservadoras
2. **Refine as instruções**: Seja mais específico sobre o que deve ser ignorado
3. **Configure ignore_files**: Adicione mais padrões de arquivos para ignorar

## Benefícios

- **Feedback imediato**: Receba comentários antes mesmo de outros desenvolvedores
- **Consistência**: Padrões uniformes de revisão
- **Aprendizado**: Aprenda boas práticas através dos comentários
- **Qualidade**: Código mais limpo e seguro
- **Produtividade**: Menos tempo gasto em revisões manuais

## Limitações

- **Não substitui revisão humana**: Use como complemento, não substituição
- **Pode ter falsos positivos**: Sempre use seu julgamento
- **Depende da qualidade das instruções**: Configure adequadamente para melhores resultados
- **Custo da API**: Cada revisão consome tokens da OpenAI

## Contribuindo

Para melhorar a revisão automática:

1. **Reporte problemas**: Abra issues para falsos positivos/negativos
2. **Sugira melhorias**: Proponha novas instruções ou configurações
3. **Teste configurações**: Experimente diferentes parâmetros
4. **Documente**: Atualize esta documentação conforme necessário
