# 🎮 Versões do Terminal - Guia Completo

## 📦 Binários Disponíveis

### **Versões de PRODUÇÃO (sem DevTools)**

```bash
./build-all.sh  # Compila ambas as versões
```

| Binário | Modo | DevTools | Uso |
|---------|------|----------|-----|
| `terminal-vertical` | Kiosk (fullscreen) | ❌ Desabilitado | Produção - Totem vertical |
| `terminal-horizontal` | Kiosk (fullscreen) | ❌ Desabilitado | Produção - Monitor horizontal |

**Características:**
- ✅ Modo kiosk (tela cheia sem controles)
- ❌ F12 desabilitado
- ❌ DevTools inacessível
- ✅ Usuário não consegue fechar/minimizar
- ✅ Ideal para produção

---

### **Versões de DEBUG (com DevTools)**

```bash
./build-debug.sh  # Compila versões debug
```

| Binário | Modo | DevTools | Uso |
|---------|------|----------|-----|
| `terminal-debug-vertical` | Windowed (janela) | ✅ **Habilitado** | Debug - Totem vertical |
| `terminal-debug-horizontal` | Windowed (janela) | ✅ **Habilitado** | Debug - Monitor horizontal |

**Características:**
- ❌ Modo windowed (janela normal)
- ✅ **F12 habilitado**
- ✅ **DevTools abre automaticamente**
- ✅ Pode fechar/minimizar
- ✅ Ideal para desenvolvimento/debug

---

## 🚀 Como Usar

### **Produção (Terminal Real)**

```bash
# Vertical (Tigrinho)
./terminal-vertical

# Horizontal (Empire)
./terminal-horizontal
```

**Você verá:**
```
🔒 DevTools DESABILITADO - Modo Produção (Kiosk)
```

**Comportamento:**
- Navegador abre em **tela cheia**
- **Sem barra de endereço/controles**
- F12 não funciona
- Esc não fecha
- Alt+F4 bloqueado

---

### **Debug (Desenvolvimento)**

```bash
# Vertical com DevTools
./terminal-debug-vertical

# Horizontal com DevTools
./terminal-debug-horizontal
```

**Você verá:**
```
🐛 DevTools HABILITADO - Modo Debug
   Pressione F12 para abrir DevTools
```

**Comportamento:**
- Navegador abre em **janela maximizada**
- **DevTools abre automaticamente** (painel lateral)
- **F12 funciona** para abrir/fechar DevTools
- Pode fechar normalmente (X)
- Pode minimizar

---

## 🔍 Debug do JavaScript

Quando usar a versão debug:

### **1. Console do Navegador**
- Logs do site: `console.log()`, `console.error()`
- Erros JavaScript
- Avisos e warnings

### **2. Network Tab**
- Todas as requisições HTTP
- POST para `/api/hooks/pnr/deposit_cash`
- WebSocket (Pusher)
- Tempo de resposta

### **3. Elements Tab**
- Inspecionar DOM
- Ver iframes dos jogos
- Modificar CSS ao vivo

### **4. Sources Tab**
- Ver código JavaScript do site
- Colocar breakpoints
- Debug passo a passo

---

## 📊 Comparação

| Feature | Produção | Debug |
|---------|----------|-------|
| **Modo** | Kiosk (fullscreen) | Windowed |
| **DevTools** | ❌ Bloqueado | ✅ Ativo |
| **F12** | ❌ Não funciona | ✅ Funciona |
| **Console logs** | ❌ Invisível | ✅ Visível |
| **Network monitor** | ❌ Não | ✅ Sim |
| **Fechar janela** | ❌ Bloqueado | ✅ Permitido |
| **Uso** | Terminais reais | Desenvolvimento |

---

## 🛠️ Builds

### **Build de Produção**
```bash
./build-all.sh
# Cria: terminal-vertical + terminal-horizontal
```

### **Build de Debug**
```bash
./build-debug.sh
# Cria: terminal-debug-vertical + terminal-debug-horizontal
```

### **Build Manual**
```bash
# Edite .env e mude ENABLE_DEVTOOLS
nano .env

# Produção
ENABLE_DEVTOOLS=false
go build -o terminal main.go

# Debug
ENABLE_DEVTOOLS=true
go build -o terminal-debug main.go
```

---

## 💡 Quando Usar Cada Versão

### **Use PRODUÇÃO quando:**
- ✅ Implantar em terminal real
- ✅ Impedir acesso do usuário aos controles
- ✅ Garantir que navegador fique em fullscreen
- ✅ Ambiente de produção/cliente

### **Use DEBUG quando:**
- ✅ Desenvolver/testar novas features
- ✅ Debugar problemas de JavaScript
- ✅ Verificar requisições de rede
- ✅ Analisar comportamento do site
- ✅ Encontrar erros no console

---

## 🔧 Configuração no .env

```env
# Debug Configuration
ENABLE_DEVTOOLS=false  # Produção
ENABLE_DEVTOOLS=true   # Debug
```

**Importante:** O valor é embarcado no binário durante a compilação!

Para mudar o modo, você precisa **recompilar**:
```bash
# Mudar .env
nano .env

# Recompilar
go build -o terminal main.go
```

---

## 📝 Exemplos de Uso

### **Exemplo 1: Debug de erro de saldo**
```bash
# 1. Compilar versão debug
./build-debug.sh

# 2. Executar
./terminal-debug-vertical

# 3. Abrir DevTools (F12 ou abre automático)

# 4. Ir para Console tab

# 5. Inserir nota (tecla 7)

# 6. Ver no console:
#    - Requisição POST enviada
#    - Resposta do servidor
#    - Qualquer erro JavaScript
```

### **Exemplo 2: Verificar comunicação WebSocket**
```bash
# 1. Executar versão debug
./terminal-debug-vertical

# 2. DevTools → Network tab → WS (WebSocket)

# 3. Ver conexão Pusher

# 4. Inserir nota

# 5. Ver mensagem 'deposit_done' no WebSocket
```

---

## ⚙️ Scripts de Build

```
build-all.sh     → Produção (vertical + horizontal)
build-debug.sh   → Debug (vertical + horizontal)
```

Ambos fazem:
1. Backup do .env
2. Modificam ENABLE_DEVTOOLS
3. Compilam versões
4. Restauram .env original

---

## 🎯 Resumo Rápido

```bash
# PRODUÇÃO (terminal real)
./build-all.sh
./terminal-vertical          # Kiosk, sem F12

# DEBUG (desenvolvimento)
./build-debug.sh
./terminal-debug-vertical    # Windowed, com F12
```

---

**🎮 Use a versão certa para cada situação!**
