# Deploy automático (GHCR + SSH)

Este guia descreve o fluxo de deploy do my-anime-api para um VPS, usando o workflow `.github/workflows/deploy.yml`.

## Visão geral do fluxo

```
push na main
  └─ build-and-push  → builda a imagem Docker e faz push para ghcr.io/<owner>/<repo>:latest
                       e ghcr.io/<owner>/<repo>:sha-<short>
  └─ deploy          → (apenas se vars.DEPLOY_ENABLED == "true")
                       SSHa no VPS, copia docker-compose.prod.yml + prometheus/ + grafana/
                       faz docker login no GHCR com o GITHUB_TOKEN, pull, up -d e prune
```

O job de build **sempre** roda. O job de deploy fica **dormente** até você setar a variável `DEPLOY_ENABLED=true` no GitHub, então é seguro fazer merge na main mesmo antes do VPS existir.

---

## 1. Preparar o VPS

Pré-requisitos no servidor:

- **Docker Engine** + **Docker Compose v2** instalados (`docker compose version` deve responder).
- Um usuário não-root (ex. `deploy`) com permissão para rodar `docker` (membro do grupo `docker`).
- Diretório onde o stack vai viver — convenção: `/home/deploy/myanimeapi` ou `/opt/myanimeapi`.
- Porta 22 (SSH) liberada para o IP do GitHub Actions OU, mais simples, qualquer IP (use chave forte).
- Portas a expor (ajuste o firewall conforme sua estratégia de reverse-proxy):
  - `8080` — API
  - `3000` — Grafana
  - `9090` — Prometheus (recomendado fechar para a internet e acessar via tunnel)

Crie o diretório de deploy e o arquivo `.env` com os segredos da aplicação **e** os parâmetros que o compose substitui:

```bash
sudo mkdir -p /opt/myanimeapi
sudo chown deploy:deploy /opt/myanimeapi
cd /opt/myanimeapi
# Crie .env com:
#   - Credenciais da app: JWT_SECRET_KEY, ALLOWED_ORIGINS, ADMIN_*, AWS_*, SMTP_* etc.
#     (veja .env.production.example no repo)
#   - Postgres: DB_USER, DB_PASSWORD, DB_NAME
#   - Observabilidade: GRAFANA_ADMIN_PASSWORD (obrigatório), GRAFANA_ADMIN_USER (opcional)
#   - Opcional: API_PORT, PROMETHEUS_PORT, GRAFANA_PORT
```

> O arquivo `.env` **não** é copiado pelo workflow — ele fica permanentemente no VPS e nunca é commitado. Trate-o como segredo de longa duração. É a única fonte de configuração: o compose lê dele tanto para substituição de variáveis (`${DB_USER}` etc.) quanto para popular o ambiente do container da API via `env_file`.

---

## 2. Gerar o par de chaves SSH para o deploy

No seu desktop (ou em qualquer máquina segura), gere um par dedicado para o GitHub Actions — **não reutilize sua chave pessoal**:

```bash
ssh-keygen -t ed25519 -C "github-actions-myanimeapi" -f ~/.ssh/myanimeapi_deploy -N ""
```

Resultado:

- `~/.ssh/myanimeapi_deploy` → chave **privada** (vai para o secret `SSH_KEY`)
- `~/.ssh/myanimeapi_deploy.pub` → chave **pública** (vai para o VPS)

No VPS, autorize a chave pública para o usuário de deploy:

```bash
# rodando como root ou via sudo
mkdir -p /home/deploy/.ssh && chmod 700 /home/deploy/.ssh
echo 'COLE-O-CONTEÚDO-DA-myanimeapi_deploy.pub-AQUI' >> /home/deploy/.ssh/authorized_keys
chmod 600 /home/deploy/.ssh/authorized_keys
chown -R deploy:deploy /home/deploy/.ssh
```

Teste o login do seu desktop:

```bash
ssh -i ~/.ssh/myanimeapi_deploy deploy@SEU_VPS_IP
```

---

## 3. Cadastrar secrets e variáveis no GitHub

Vá em **Settings → Secrets and variables → Actions** no repositório.

### Secrets (Repository secrets)

| Nome          | Obrigatório | Valor                                                                                  |
|---------------|-------------|----------------------------------------------------------------------------------------|
| `SSH_HOST`    | sim         | IP ou hostname do VPS (ex. `203.0.113.42` ou `vps.meudominio.com`)                     |
| `SSH_USER`    | sim         | Usuário SSH no VPS (ex. `deploy`)                                                       |
| `SSH_KEY`     | sim         | Conteúdo **completo** do arquivo de chave privada (`~/.ssh/myanimeapi_deploy`), incluindo as linhas `-----BEGIN ...-----` e `-----END ...-----` |
| `SSH_PORT`    | não         | Porta SSH se diferente de 22 (ex. `2222`)                                              |
| `DEPLOY_PATH` | sim         | Caminho absoluto no VPS onde o stack vive (ex. `/opt/myanimeapi`)                       |

Observações:

- `GITHUB_TOKEN` é fornecido automaticamente pelo Actions — você **não** precisa cadastrar. Ele é usado para autenticar tanto no push da imagem para o GHCR quanto no pull no servidor durante o deploy.
- A chave privada (`SSH_KEY`) precisa ser a chave **inteira**, sem alterar quebras de linha. Cole direto do arquivo.

### Variables (Repository variables)

| Nome              | Valor                                                                                  |
|-------------------|----------------------------------------------------------------------------------------|
| `DEPLOY_ENABLED`  | `true` para ativar o job de deploy. Enquanto não estiver setado, só o build+push roda. |

> Use uma **variable**, não um secret, porque ela precisa ser legível na condição `if:` do job. Não é um valor sensível.

---

## 4. Configuração do GHCR (uma vez só)

Após o primeiro push bem-sucedido, vá em **github.com/<seu-usuário>?tab=packages**, abra o package `my-anime-api` (ou o nome do repo) e:

- Em **Package settings → Manage Actions access**, garanta que o repo `my-anime-api` tem acesso `Write`.
- Em **Danger Zone → Change package visibility**, decida se a imagem é pública ou privada. Se privada, o servidor precisa logar no GHCR a cada pull — o workflow já cuida disso durante o deploy.

---

## 5. Primeiro deploy

1. Confirme que `.env` existe em `$DEPLOY_PATH` no VPS.
2. Faça push na `main` (ou rode o workflow manualmente em **Actions → Build & Deploy → Run workflow**).
3. Acompanhe os dois jobs no Actions. Se `DEPLOY_ENABLED` ainda não estiver `true`, só o `build-and-push` roda.
4. Quando estiver pronto, set `DEPLOY_ENABLED=true` e re-rode o workflow.

Após o `up -d`, valide:

```bash
curl http://SEU_VPS_IP:8080/v1/health        # {"status":"success",...}
curl http://SEU_VPS_IP:8080/metrics | head    # métricas Prometheus
# Grafana em http://SEU_VPS_IP:3000 (usuário/senha conforme .env.production)
```

---

## 6. Operação manual no VPS

Restart sem fazer pull (usa a imagem já em disco):

```bash
cd /opt/myanimeapi
export IMAGE_REF=ghcr.io/<owner>/<repo>:latest
docker compose -f docker-compose.prod.yml up -d
```

Pull manual de uma nova tag (precisa estar logado no GHCR — para repos privados):

```bash
echo "<seu-PAT-com-read:packages>" | docker login ghcr.io -u <seu-user> --password-stdin
export IMAGE_REF=ghcr.io/<owner>/<repo>:sha-abc1234
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

Logs:

```bash
docker compose -f docker-compose.prod.yml logs -f api
docker compose -f docker-compose.prod.yml logs -f prometheus grafana
```

---

## 7. Resumo dos secrets/variables

Para revisar rapidamente, o que você precisa cadastrar quando o VPS estiver no ar:

**Secrets:** `SSH_HOST`, `SSH_USER`, `SSH_KEY`, `DEPLOY_PATH` (e `SSH_PORT` se não for 22).
**Variable:** `DEPLOY_ENABLED=true`.

E no VPS: Docker + Compose v2 instalados, usuário de deploy no grupo `docker`, diretório `$DEPLOY_PATH` existente com `.env.production` dentro.
