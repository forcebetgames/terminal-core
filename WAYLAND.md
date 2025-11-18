# Suporte Wayland/X11

## 🎉 O programa agora funciona em AMBOS os ambientes!

Este terminal detecta automaticamente se você está usando **X11** ou **Wayland** e utiliza o backend apropriado para captura de teclado.

---

## 🖥️ Detecção Automática

Ao iniciar, você verá uma mensagem indicando qual ambiente foi detectado:

### X11
```
🖥️  Display Server detectado: x11
🪟 Usando modo X11 (gohook)
🎧 Iniciando escuta de eventos gohook (X11)...
```

### Wayland
```
🖥️  Display Server detectado: wayland
🌊 Usando modo WAYLAND (evdev)
⚠️  IMPORTANTE: Execute com sudo se houver erros de permissão!
✅ Teclado detectado: AT Translated Set 2 keyboard (/dev/input/event3)
🎧 Iniciando escuta de eventos evdev (Wayland)...
```

---

## 🚀 Como Executar

### No X11 (padrão)
```bash
./terminal-horizontal
# ou
./terminal-vertical
```

### No Wayland
```bash
sudo ./terminal-horizontal
# ou
sudo ./terminal-vertical
```

**⚠️ IMPORTANTE no Wayland**: O acesso a `/dev/input/event*` requer permissões de root, então você **deve executar com `sudo`**.

---

## 🔍 Como Verificar Qual Ambiente Você Está Usando?

```bash
echo $XDG_SESSION_TYPE
```

**Resultado:**
- `x11` → Você está usando X11
- `wayland` → Você está usando Wayland

---

## ⚙️ Diferenças Técnicas

| Aspecto | X11 (gohook) | Wayland (evdev) |
|---------|--------------|-----------------|
| **Biblioteca** | robotn/gohook | gvalkov/golang-evdev |
| **Permissões** | Usuário normal | **Requer sudo** |
| **Detecção de teclas** | Via X11 APIs | Leitura direta do kernel |
| **Compatibilidade** | X11 apenas | Funciona em qualquer ambiente |
| **Performance** | Boa | Excelente |

---

## 🐛 Troubleshooting

### Problema: "nenhum teclado encontrado" no Wayland
**Solução**: Execute com `sudo`
```bash
sudo ./terminal-horizontal
```

### Problema: "XkbGetKeyboard failed" no X11
**Solução**: Instale as dependências do X11:
```bash
sudo pacman -S libx11 libxtst xorg-server-devel libxcb libxkbcommon libxkbcommon-x11
```

### Problema: Teclas não funcionam em nenhum ambiente
**Solução**: Verifique os logs de debug. Você deve ver mensagens como:
```
🔍 DEBUG - Tecla capturada:
   Tecla: 7
   Rawcode: 8
   Valor: R$ 100
✅ VÁLIDO - Inserindo R$ 100...
```

Se não vir essas mensagens, o hook não está capturando eventos.

---

## 📦 Mapeamento de Teclas

O sistema detecta as seguintes teclas em **ambos** os ambientes:

| Tecla | Valor | evdev Keycode | Descrição |
|-------|-------|---------------|-----------|
| `2` ou `Numpad 2` | R$ 2 | 3 ou 80 | Dois reais |
| `3` ou `Numpad 3` | R$ 5 | 4 ou 81 | Cinco reais |
| `4` ou `Numpad 4` | R$ 10 | 5 ou 82 | Dez reais |
| `5` ou `Numpad 5` | R$ 20 | 6 ou 83 | Vinte reais |
| `6` ou `Numpad 6` | R$ 50 | 7 ou 84 | Cinquenta reais |
| `7` ou `Numpad 7` | R$ 100 | 8 ou 85 | Cem reais |

---

## 🔐 Permissões no Wayland (Opcional)

Se você não quiser executar com `sudo` toda vez, pode adicionar seu usuário ao grupo `input`:

```bash
sudo usermod -aG input $USER
```

Depois, faça **logout e login** novamente para que as permissões tenham efeito.

**⚠️ Aviso de Segurança**: Isso permite que qualquer programa executado pelo seu usuário leia eventos de teclado globalmente. Use com cautela.

---

## 🎯 Arquitetura do Código

### Abstração `InputHandler`
```
keyboard.InputHandler (interface)
    ├── RealInputHandler (X11 via gohook)
    └── EvdevInputHandler (Wayland via evdev)
```

### Detecção Automática
```go
// domain/keyboard/evdev-handler.go
func NewInputHandler() (InputHandler, error) {
    displayServer := DetectDisplayServer()

    if displayServer == "wayland" {
        return NewEvdevInputHandler()
    }

    return NewRealInputHandler(), nil
}
```

### Arquivos Relevantes
- `domain/keyboard/input-handler.go` - Backend X11 (gohook)
- `domain/keyboard/evdev-handler.go` - Backend Wayland (evdev)
- `domain/keyboard/event.go` - Interface comum
- `domain/payment.go` - Sistema de pagamento (usa a abstração)
- `main.go` - Inicialização e detecção automática

---

## ✨ Vantagens do Suporte Híbrido

1. **Compatibilidade Universal** - Funciona em qualquer distribuição Linux moderna
2. **Zero Configuração** - Detecção automática do ambiente
3. **Performance Otimizada** - Cada backend é otimizado para seu ambiente
4. **Fallback Inteligente** - Se um backend falhar, você pode mudar para o outro

---

## 📝 Logs de Debug

O sistema fornece logs detalhados para facilitar o troubleshooting:

### Eventos Válidos
```
====================================
🔍 DEBUG - Tecla capturada:
   Tecla: 7
   Rawcode: 8
   Valor: R$ 100
====================================
✅ VÁLIDO - Inserindo R$ 100...
====================================
🚀 CALLBACK EXECUTADO - Valor: R$ 100
====================================
```

### Eventos Rejeitados (Debounce)
```
====================================
🔍 DEBUG - Tecla capturada:
   Tecla: 7
   Rawcode: 8
   Valor: R$ 100
   ⏱️  REJEITADO: Debounce ativo (última tecla há 50ms)
====================================
```

---

## 🔄 Como Mudar entre X11 e Wayland

Se você quiser testar em ambos os ambientes, faça logout e na tela de login:

1. Clique no ícone de engrenagem ⚙️ (geralmente no canto inferior direito)
2. Selecione:
   - **GNOME** (Wayland) → Usa Wayland
   - **GNOME on Xorg** → Usa X11

Ou permanentemente via `/etc/gdm/custom.conf`:
```ini
[daemon]
# Descomente para forçar X11
#WaylandEnable=false
```
