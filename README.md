# GOWSAA
Ejemplo en lenguaje GO de conexión al servicio de autenticación de AFIP (webservice wsaa)

#### Requerimientos
1. Variable de ambiente/entorno (AfipP12) con nombre y ubicación del archivo .p12 (certificado y privateKey RSA). Para entorno Windows utilizar doble barra como separador de directorios, ej. c:\\archivo.p12
2. Variable de ambiente/entorno (AfipP12Pass) con contraseña del archivo .p12 del punto 1

#### Ejecución
1. Descargar los fuentes
2. Ejecutar alguna de estas opciones: 
    * go run .\cmd\main.go
    * go build -o gowsaa.exe .\cmd\main.go



#### Créditos
Para la conexión soap y para la generación del archivo wsaa.go se utilizó https://github.com/hooklift/gowsdl/


