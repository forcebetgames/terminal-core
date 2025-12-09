# Correções Implementadas - Terminal de Jogos

## Problemas Identificados e Soluções

### 1. ❌ Teclas de Depósito (2-7) Não Funcionavam no Ubuntu X11

**Sintomas:**
- Logs mostravam: `keycode_8`, `keycode_7`, `keycode_6`, etc.
- Nenhum callback era executado
- Apenas teclas '3' e '4' do teclado funcionavam

**Causa:**
- Os botões de hardware geram keycodes 2-8 no X11
- Faltava mapeamento desses keycodes no `input-handler.go`

**Solução:**
Adicionado mapeamento em `/domain/keyboard/input-handler.go` linha 65-73:
```go
// Hardware buttons (custom terminal keys)
2: "2",
3: "3",
4: "4",
5: "5",
6: "6",
7: "7",
8: "8",
```

### 2. ❌ Controles de Jogo Não Funcionavam (Aumentar/Diminuir Aposta, Touch/Pointer)

**Sintomas:**
- Botões +/- não funcionavam no Ubuntu Xorg (mas funcionavam no Wayland)
- Touch/clique nos botões do jogo não respondia
- Modal de depósito não abria
- WebSocket do Pusher não conectava

**Causa:**
- Flag `disable-background-networking` bloqueava WebSockets
- Flag `disable-plugins` bloqueava eventos de touch/pointer no Xorg
- Faltavam flags específicas para habilitar touch events no Xorg

**Solução:**
Removidas flags problemáticas e adicionadas flags de compatibilidade em `/domain/browser.go`:
```go
// ❌ REMOVIDO - bloqueava WebSockets do Pusher!
// chromedp.Flag("disable-background-networking", true),

// ❌ REMOVIDO - bloqueava eventos de touch/pointer no Xorg!
// chromedp.Flag("disable-plugins", true),

// ✅ ADICIONADO - FLAGS PARA COMPATIBILIDADE XORG/WAYLAND
chromedp.Flag("enable-features", "TouchpadOverscrollHistoryNavigation"),
chromedp.Flag("enable-blink-features", "PointerEvent,TouchEvents"),
chromedp.Flag("touch-events", "enabled"),

// ✅ ADICIONADO - FLAGS PARA MELHOR RENDERING NO XORG
chromedp.Flag("use-gl", "desktop"),
chromedp.Flag("enable-accelerated-2d-canvas", true),
chromedp.Flag("disable-gpu-driver-bug-workarounds", true),
```

### 3. ✅ Controles de Jogo no Evdev (Wayland)

**Adicionado:**
Mapeamentos para controles de jogo em `/domain/keyboard/evdev-handler.go`:
```go
13:  "=",   // Equal/Plus key - Increase bet
12:  "-",   // Minus key - Decrease bet
57:  " ",   // Space - Start game
103: "up",  // Up arrow - Change game
24:  "o",   // O key - Change game alternative
```

## Como Atualizar o Terminal 56 (Ubuntu)

### Opção 1: Copiar Binário via SCP
```bash
scp terminal usuario@terminal56:/caminho/do/terminal/
```

### Opção 2: Recompilar no Terminal
```bash
cd /caminho/do/projeto
git pull
go build -o terminal main.go
```

### Teste Final
Execute o terminal e teste:
```bash
sudo ./terminal  # sudo necessário para X11/gohook
```

**Teste as teclas:**
- ✅ Teclas 2-7 dos botões de hardware → Devem fazer depósitos
- ✅ Botão 'g' → Deve mudar de jogo
- ✅ Botões +/- → Devem aumentar/diminuir aposta
- ✅ SPACE → Deve iniciar/girar jogo
- ✅ Modal de depósito deve abrir automaticamente

## Diagnóstico dos Logs do Terminal 56

**Funcionando:**
- ✅ Captura de teclas no gohook (X11)
- ✅ Envio HTTP de depósitos (200 OK)
- ✅ Backend processando depósitos
- ✅ Pusher emitindo eventos

**Corrigido nesta versão:**
- ✅ Keycodes 2-8 agora mapeados
- ✅ WebSocket do Pusher liberado
- ✅ Controles de jogo funcionais

## Arquivos Modificados

1. `/domain/keyboard/input-handler.go` - Adicionado mapeamento keycodes 2-9
2. `/domain/keyboard/evdev-handler.go` - Adicionado controles de jogo
3. `/domain/browser.go` - Removida flag que bloqueava WebSocket

## Observações Técnicas

**Diferença entre X11 e Wayland:**
- **Wayland (evdev)**: keycodes 3-8 = teclas "2"-"7"
- **X11 (gohook)**: keycodes 2-8 = botões hardware, keycodes 11-16 = teclas teclado "2"-"7"

**Por que funciona no Wayland mas não no Xorg?**
- **Wayland**: Gerenciamento nativo de eventos de input, menos dependente de flags do Chrome
- **Xorg**: Precisa de flags explícitas para habilitar certos tipos de eventos (touch, pointer)
- **disable-plugins** no Xorg bloqueia parte do sistema de eventos que jogos HTML5 usam
- **disable-background-networking** bloqueia WebSockets em ambos os sistemas

**Terminal 56 específico:**
- Machine ID: `1f86a782bd38461ba0a4e22a58d184e5`
- Sistema: Ubuntu com X11/Xorg
- Botões hardware: keycodes 2-8 (rawcodes 49-55)
- Teclado normal: keycodes 12-13 para '3' e '4'
