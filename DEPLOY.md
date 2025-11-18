# 🚀 Guia de Deploy - Binário Único Universal

## 📦 Arquivo de Deploy

Agora você tem **UM ÚNICO BINÁRIO** que funciona em qualquer máquina:

```
terminal  (binário universal - funciona em horizontal E vertical)
```

## ✨ Detecção Automática

O binário detecta **automaticamente** a orientação da tela:

1. **Detecção automática via xrandr (X11/Wayland)**
   - Lê a resolução atual da tela
   - Se largura > altura → HORIZONTAL (1920x1080)
   - Se altura > largura → VERTICAL (1080x1920)

2. **Override manual via linha de comando**
   ```bash
   ./terminal horizontal  # Força modo horizontal
   ./terminal vertical    # Força modo vertical

   # Aliases aceitos:
   ./terminal landscape
   ./terminal portrait
   ./terminal -h
   ./terminal -v
   ```

3. **Override via variável de ambiente**
   ```bash
   SCREEN_ORIENTATION=landscape ./terminal
   SCREEN_ORIENTATION=portrait ./terminal
   ```

## 🎯 Deploy Simplificado

### Deploy em TODAS as máquinas (mesmo binário!)

```bash
# 1. Compilar o binário
go build -o terminal main.go

# 2. Copiar para QUALQUER máquina (vertical ou horizontal)
scp terminal usuario@maquina:/opt/terminal/

# 3. Executar (detecta automaticamente!)
ssh usuario@maquina
cd /opt/terminal
sudo ./terminal  # Detecta se é horizontal ou vertical sozinho!
```

### Deploy com Wayland

```bash
# Execute com sudo (necessário para evdev)
sudo ./terminal
```

## 📊 Exemplos de Uso

### Máquina com tela horizontal (1920x1080)
```bash
./terminal
# 🔍 Resolução detectada: 1920x1080 → HORIZONTAL
# 📐 Modo: HORIZONTAL (Landscape) - 1920x1080
# 🎯 Jogo: Empire 🏛️
```

### Máquina com tela vertical (1080x1920)
```bash
./terminal
# 🔍 Resolução detectada: 1080x1920 → VERTICAL
# 📐 Modo: VERTICAL (Portrait) - 1080x1920
# 🎯 Jogo: Tigrinho 🐯
```

### Forçar modo específico (ignorar detecção)
```bash
# Forçar horizontal mesmo em tela vertical
./terminal horizontal

# Forçar vertical mesmo em tela horizontal
./terminal vertical
```

## 🎮 Vantagens

1. **Deploy Único**: Um arquivo só para todas as máquinas
2. **Zero Configuração**: Detecta automaticamente
3. **Override Simples**: Pode forçar manualmente se necessário
4. **Menor Manutenção**: Não precisa lembrar qual versão enviar
5. **Menos Espaço**: Um binário em vez de dois

## 🔧 Compilação

### Compilar binário universal
```bash
go build -o terminal main.go
```

### Compilar versões antigas (se necessário)
```bash
# Ainda funciona se quiser binários separados:
SCREEN_ORIENTATION=portrait go build -o terminal-vertical main.go
SCREEN_ORIENTATION=landscape go build -o terminal-horizontal main.go
```

## 🐛 Troubleshooting

### Problema: Detectou orientação errada

**Solução 1**: Force manualmente
```bash
./terminal horizontal
```

**Solução 2**: Use variável de ambiente
```bash
SCREEN_ORIENTATION=landscape ./terminal
```

### Problema: "Permission denied" no Wayland

**Solução**: Execute com sudo
```bash
sudo ./terminal
```

### Problema: Não detectou resolução

Verifique se xrandr está instalado:
```bash
# Instalar xrandr (Arch Linux)
sudo pacman -S xorg-xrandr

# Testar detecção
xrandr --current
```

## 📋 Checklist de Deploy

- [ ] Compilar: `go build -o terminal main.go`
- [ ] Testar localmente: `./terminal`
- [ ] Verificar que detectou corretamente
- [ ] Copiar para máquina remota: `scp terminal user@host:/opt/terminal/`
- [ ] Executar remotamente: `ssh user@host 'cd /opt/terminal && sudo ./terminal'`
- [ ] Validar orientação e jogo corretos

## 💡 Dicas

1. **Sempre use o binário único** - mais simples e menos chance de erro
2. **Teste localmente primeiro** - verifique se a detecção funciona
3. **Use sudo no Wayland** - necessário para captura de teclado
4. **Override só quando necessário** - deixe a detecção automática funcionar

## 🎯 Comparação: Antes vs Agora

### ❌ Antes (2 binários)
```bash
# Tinha que lembrar qual enviar:
scp terminal-vertical usuario@maquina-vertical:/opt/
scp terminal-horizontal usuario@maquina-horizontal:/opt/

# Errar = jogo errado na tela!
```

### ✅ Agora (1 binário)
```bash
# Sempre o mesmo comando:
scp terminal usuario@qualquer-maquina:/opt/

# Detecta automaticamente = sem erros!
```

---

**🎮 Deploy simplificado com detecção automática - ForceBet Games**
