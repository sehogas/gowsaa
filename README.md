# GOWSAA
Ejemplo de conexión al servicio de autenticación de AFIP (webservice wsaa) en lenguaje GO.

#### Requisitos previos
1. Generar el archivo con extensión .p12 contenedor del certificado y clave privada RSA. [Ver documentación de la AFIP](https://www.afip.gob.ar/ws/documentacion/wsaa.asp).
2. Variable de ambiente/entorno **AfipP12** con ruta y nombre del archivo .p12. Para entorno Windows utilizar doble barra como separador de directorios, ej. c:\\\archivo.p12
3. Variable de ambiente/entorno **AfipP12Pass** con contraseña del archivo .p12

#### Ejecución
1. Descargar los fuentes

``git clone https://github.com/sehogas/gowsaa.git``

2. Ejecutar alguna de estas opciones:

  ``go run .\cmd\gowsaa.go``

  ``go build .\cmd\gowsaa.go``


---
#### Créditos
Para la conexión soap y la generación del archivo wsaa.go se utilizó https://github.com/hooklift/gowsdl/