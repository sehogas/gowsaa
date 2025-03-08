# GOWSAA
Ejemplo de conexión al servicio de autenticación de AFIP (webservice wsaa) en lenguaje GO.

#### Requisitos previos
1. Generar clave privada RSA, crear solicitud de certificado y obtener en AFIP el certificado. 
2. Configurar variables de ambiente con los datos del paso 1.

#### Ejecución
1. Descargar los fuentes
``git clone https://github.com/sehogas/gowsaa.git``

2. Ejecutar alguna de estas opciones:

  ``go run .\cmd\gowssa.go``

  ``go build .\cmd\gowsaa.go``


#### Pasos para generar clave privada RSA y Certificado AFIP
  Documentación: https://www.afip.gob.ar/ws/WSASS/html/generarcsr.html

    openssl genrsa -out MiClavePrivada 2048

    openssl req -new -key MiClavePrivada -subj "/C=AR/O=XXXX/CN=YYYY/serialNumber=CUIT 20999999992" -out misolicitud.csr

  Crear el archivo "certificado.pem" y copiar el certificado x509v2 en formato PEM generado por la página de AFIP 

#### Configurar variables de ambiente 
  PRIVATE_KEY_FILE=MiClavePrivada
  CERTIFICATE_FILE=certificado.pem

---
#### Créditos
  Para la conexión soap y la generación del archivo wsaa.go se utilizó https://github.com/hooklift/gowsdl/