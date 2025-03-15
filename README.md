# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Además, definir un script de bash `generar-compose.sh` que permita crear una definición de DockerCompose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```
#### Resolución:

Se creo un script llamado "generar-compose-sh" el cual llama a un script en python llamado "mi-generador.py" que genera efectivamente el docker-compose con la cantidad de clientes solicitada.

#### Ejecución:
```
. generar-compose.sh docker-compose-dev.yaml <nro clientes>
make docker-compose-up
docker ps
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°2:

Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (config.ini y config.yaml, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: docker volumes).

#### Resolución:

Modifique "mi-generador.py" con el objetivo de que cree volumes para cada tipo de servicio (cliente y servidor) sobre los archivos de configuración.

#### Ejecución:
```
. generar-compose.sh docker-compose-dev.yaml <nro clientes>
make docker-compose-up
docker ps
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `

#### Resolución:

Creo el archivo "docker-compose-validar-dev.yaml" el cual levanta un cliente netcat en la misma red que server, asumiendo que el contenedor de "server" ya esta corriendo. Este contenedor, "netcat-cli", corre un script llamado "do-validate-echo-server" cuyo log es parseado por el script "validar-echo-server.sh" el cual imprime el resultado de el mensaje al servidor.

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 0
make docker-compose-up
. validar-echo-server.sh
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

#### Resolución:

En el cliente se crea un canal en el "main" al cual se notifican las señales las cuales son tomadas del canal en el loop principal de cliente en "client.go" para finalizar dicho loop correctamente.
En el servidor se crea una función que introduce un callback para el manejo de señales llamado "__handle_a_signal" el cual efectúa el log del evento cierra los socket's de cada uno de los clientes para luego finalizar el proceso servidor.

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 1
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```
NOTA: Para verificar que el cierre del cliente funciona correctamente subir el loop.amount en el archivo de configuración del mismo retornando un "code 0" en el log.

***

#### Test:

```
REPO_PATH=/home/putaro/Workspace/tp0 pytest -s
```
***