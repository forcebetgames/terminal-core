# Fluxo
- Iniciar as variáveis de ambiente
- iniciar o banco de dados
- identificar o terminal
    - Validar o status do terminal
    - Se o não achar, não iniciar o programa
    - Atualizar ip publico
- habilitar websocket
- iniciar o browser


# Como descobrir o machine id

### Windows
wmic csproduct get UUID

### Ubuntu Linux
cat /etc/machine-id

### Macos
ioreg -rd1 -c IOPlatformExpertDevice | grep IOPlatformUUID