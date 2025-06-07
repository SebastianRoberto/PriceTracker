@echo off
echo Configurando reglas de firewall para la aplicación...

:: Eliminar reglas existentes si existen (para evitar duplicados)
netsh advfirewall firewall delete rule name="Go App - HTTP (8080)" > nul 2>&1

:: Agregar regla para el puerto HTTP (8080)
netsh advfirewall firewall add rule name="Go App - HTTP (8080)" dir=in action=allow protocol=TCP localport=8080

:: Agregar regla para el proceso de Go
netsh advfirewall firewall add rule name="Go Application" dir=in action=allow program="%USERPROFILE%\go\bin\go.exe" enable=yes

echo Configuración del firewall completada.
echo.
echo Ahora puedes ejecutar la aplicación sin alertas del firewall.
echo.
pause 