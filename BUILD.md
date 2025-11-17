# 🏗️ Guia de Build - Terminal de Jogos

## 🚀 Build Rápido (Ambas Versões)

Para compilar **VERTICAL** e **HORIZONTAL** de uma vez:

```bash
./build-all.sh
```

Isso vai criar:
- `terminal-vertical` (1080x1920 - Vídeo do Tigrinho 🐯)
- `terminal-horizontal` (1920x1080 - Vídeo do Empire 🏛️)

---

## 🔨 Build Manual (Uma Versão por Vez)

### Versão VERTICAL (Tigrinho)

```bash
# 1. Edite o .env
nano .env
# Certifique-se que tem: SCREEN_ORIENTATION=portrait

# 2. Compile
go build -o terminal-vertical main.go
```

### Versão HORIZONTAL (Empire)

```bash
# 1. Edite o .env
nano .env
# Certifique-se que tem: SCREEN_ORIENTATION=landscape

# 2. Compile
go build -o terminal-horizontal main.go
```

---

## 📱 Como Usar os Binários

### Testar Localmente

```bash
# Versão Vertical (Tigrinho)
./terminal-vertical

# Versão Horizontal (Empire)
./terminal-horizontal
```

### Distribuir para Máquinas

**Terminal com tela vertical:**
```bash
scp terminal-vertical usuario@maquina-vertical:/opt/terminal/
ssh usuario@maquina-vertical
cd /opt/terminal
./terminal-vertical
```

**Terminal com tela horizontal:**
```bash
scp terminal-horizontal usuario@maquina-horizontal:/opt/terminal/
ssh usuario@maquina-horizontal
cd /opt/terminal
./terminal-horizontal
```

---

## 🎯 Configurações de Tela

### Arquivo `.env`

```env
# Vertical (Tigrinho) - 1080x1920
SCREEN_ORIENTATION=portrait

# Horizontal (Empire) - 1920x1080
SCREEN_ORIENTATION=landscape

# Custom (especifique dimensões exatas)
# SCREEN_WIDTH=1200
# SCREEN_HEIGHT=1600
```

---

## 📊 Diferenças Entre as Versões

| Característica | Vertical | Horizontal |
|---------------|----------|------------|
| Resolução | 1080x1920 | 1920x1080 |
| Vídeo de Fundo | Tigrinho 🐯 | Empire 🏛️ |
| Uso | Totem vertical | Monitor horizontal |
| Tamanho | ~16 MB | ~16 MB |

---

## 🔧 Builds Avançados

### Build para Windows

```bash
GOOS=windows GOARCH=amd64 go build -o terminal-vertical.exe main.go
GOOS=windows GOARCH=amd64 go build -o terminal-horizontal.exe main.go
```

### Build para macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o terminal-vertical-mac main.go
GOOS=darwin GOARCH=amd64 go build -o terminal-horizontal-mac main.go
```

### Build com Compressão (UPX)

```bash
# Compile normalmente
./build-all.sh

# Comprima os binários (reduz ~60%)
upx --best --lzma terminal-vertical
upx --best --lzma terminal-horizontal
```

---

## 🐛 Troubleshooting

### Erro: "Permission denied"

```bash
chmod +x build-all.sh
chmod +x terminal-vertical
chmod +x terminal-horizontal
```

### Erro: "bad interpreter"

```bash
# Corrigir line endings (Windows → Linux)
sed -i 's/\r$//' build-all.sh
```

### Recompilar após mudanças no .env

```bash
# Sempre recompile se mudar o .env
./build-all.sh
```

---

## 📦 Estrutura de Arquivos

```
terminal-core/
├── .env                    # Configurações (embarcadas no build)
├── main.go                 # Código principal
├── domain/
│   ├── browser.go         # Lógica do navegador + viewport
│   └── ...
├── build-all.sh           # Script de build automático ✨
├── terminal-vertical      # Binário vertical (portrait)
├── terminal-horizontal    # Binário horizontal (landscape)
└── terminal               # Binário padrão (vertical)
```

---

## 💡 Dicas

1. **Sempre use `./build-all.sh`** para garantir consistência
2. **Teste localmente** antes de distribuir
3. **Verifique o Machine ID** da máquina de destino
4. **Garanta que o terminal tem jogos** no banco de dados
5. **Use nomes descritivos** ao renomear os binários

---

## ✅ Checklist de Deploy

- [ ] Executar `./build-all.sh`
- [ ] Testar `terminal-vertical` localmente
- [ ] Testar `terminal-horizontal` localmente
- [ ] Verificar que ambos mostram o vídeo correto
- [ ] Obter Machine ID da máquina de destino
- [ ] Cadastrar terminal no banco
- [ ] Associar jogos ao terminal
- [ ] Copiar binário apropriado
- [ ] Executar e validar

---

## 📞 Suporte

Se algo não funcionar:

1. Verifique os logs no console
2. Execute `./check_machine_id` para ver o ID
3. Confirme que o terminal existe no banco
4. Verifique se há jogos associados
5. Teste a URL manualmente no navegador

---

**🎮 Build gerado com ❤️ para ForceBet Games**
