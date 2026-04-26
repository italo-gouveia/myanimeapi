# 🧪 Testando a Revisão Automática de Código

Este guia mostra como testar se a revisão automática de código com IA está funcionando corretamente.

## 🚀 Teste Rápido

### 1. Criar PR de Teste

1. **Crie uma branch de teste:**
   ```bash
   git checkout -b test-ai-review
   ```

2. **Faça uma mudança simples:**
   ```go
   // Adicione este código em um arquivo .go existente
   func testFunction() {
       // Este código tem alguns problemas para a IA detectar
       var password = "123456"  // Senha hardcoded
       fmt.Println("Debug info:", password)  // Log de senha
       
       // Sem validação de entrada
       userInput := getUserInput()
       executeCommand(userInput)  // Possível injeção de comando
   }
   ```

3. **Commit e push:**
   ```bash
   git add .
   git commit -m "test: adicionar código para testar revisão IA"
   git push origin test-ai-review
   ```

4. **Criar Pull Request:**
   - Vá para GitHub
   - Clique em "Compare & pull request"
   - Título: "Test: Revisão Automática de Código"
   - Descrição: "Testando se a IA detecta problemas de segurança e qualidade"

### 2. Aguardar Análise

- **Tempo esperado**: 2-5 minutos
- **Verificar**: Aba "Actions" no GitHub
- **Logs**: Clique no workflow "AI Code Review"

### 3. Verificar Comentários

A IA deve detectar:
- 🔒 **Senha hardcoded** (`password = "123456"`)
- 🔒 **Log de dados sensíveis** (`fmt.Println(password)`)
- 🔒 **Possível injeção de comando** (`executeCommand(userInput)`)
- 🔧 **Falta de validação de entrada**

## 📋 Checklist de Verificação

### ✅ Configuração Básica

- [ ] Secret `OPENAI_API_KEY` configurado
- [ ] Workflow `.github/workflows/code-review.yml` existe
- [ ] Workflow executa em PRs para branch `main`
- [ ] Permissões configuradas corretamente

### ✅ Funcionamento

- [ ] Workflow inicia automaticamente
- [ ] Não há erros nos logs
- [ ] Comentários aparecem no PR
- [ ] Comentários são relevantes e úteis

### ✅ Qualidade dos Comentários

- [ ] Comentários são específicos
- [ ] Sugestões de correção incluídas
- [ ] Emojis usados para categorização
- [ ] Foco em Go e APIs

## 🐛 Problemas Comuns

### Workflow não executa

**Sintomas:**
- Nenhuma action aparece na aba Actions
- PR não mostra comentários da IA

**Soluções:**
1. Verifique se o PR é para a branch `main`
2. Confirme que os arquivos estão na branch correta
3. Verifique se o workflow está ativo

### Erro de API Key

**Sintomas:**
- Logs mostram "Invalid API key"
- Workflow falha com erro 401

**Soluções:**
1. Verifique se `OPENAI_API_KEY` está correto
2. Confirme que a chave não expirou
3. Verifique se há créditos na conta OpenAI

### Comentários muito genéricos

**Sintomas:**
- Comentários vagos como "boa prática"
- Sem exemplos específicos

**Soluções:**
1. Ajuste as instruções em `.github/pr-agent.toml`
2. Use modelo GPT-4 em vez de GPT-3.5
3. Aumente `max_tokens`

## 🧪 Casos de Teste

### Teste 1: Problemas de Segurança

```go
// Código com problemas de segurança
func insecureHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    
    // SQL Injection vulnerável
    sql := "SELECT * FROM users WHERE name = '" + query + "'"
    db.Query(sql)
    
    // XSS vulnerável
    fmt.Fprintf(w, "<h1>Resultado: %s</h1>", query)
}
```

**Esperado**: IA deve detectar SQL injection e XSS

### Teste 2: Problemas de Performance

```go
// Código com problemas de performance
func slowFunction() {
    for i := 0; i < 1000000; i++ {
        // Operação custosa em loop
        time.Sleep(1 * time.Millisecond)
    }
}
```

**Esperado**: IA deve sugerir otimizações

### Teste 3: Problemas de Qualidade

```go
// Código com problemas de qualidade
func badFunction(x int, y int, z int) {
    var a = x + y + z
    fmt.Println(a)
}
```

**Esperado**: IA deve sugerir melhor nomenclatura e documentação

## 📊 Métricas de Sucesso

### Taxa de Detecção

- **Problemas de segurança**: >90%
- **Problemas de performance**: >80%
- **Problemas de qualidade**: >85%

### Qualidade dos Comentários

- **Específicos**: Comentários devem ser específicos
- **Acionáveis**: Sugestões devem ser implementáveis
- **Educativos**: Devem explicar o porquê

### Tempo de Resposta

- **Tempo médio**: <5 minutos
- **Tempo máximo**: <10 minutos

## 🔄 Teste Contínuo

### Automatizar Testes

1. **Criar branch de teste permanente**
2. **Script de teste automatizado**
3. **Monitoramento de métricas**

### Exemplo de Script

```bash
#!/bin/bash
# test-ai-review.sh

echo "🧪 Testando Revisão Automática de Código..."

# Criar branch de teste
git checkout -b test-ai-review-$(date +%s)

# Adicionar código de teste
cat > test_file.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
)

func insecureHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    sql := "SELECT * FROM users WHERE name = '" + query + "'"
    fmt.Fprintf(w, "<h1>Resultado: %s</h1>", query)
}
EOF

# Commit e push
git add test_file.go
git commit -m "test: adicionar código inseguro para teste"
git push origin HEAD

echo "✅ PR criado. Verifique os comentários da IA em 5 minutos."
```

## 📈 Melhorias Contínuas

### Coletar Feedback

1. **Avaliar comentários** da IA
2. **Ajustar configurações** conforme necessário
3. **Treinar equipe** sobre como usar feedback

### Refinar Configurações

1. **Ajustar instruções** baseado no feedback
2. **Otimizar custos** vs qualidade
3. **Personalizar** para o projeto específico

---

**Dica**: Mantenha um PR de teste sempre aberto para verificar se o sistema está funcionando corretamente.
