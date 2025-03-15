import sys

ERROR_MESSAGE = "No se han ingresado los argumentos correctos: <archivo de salida> <cantidad clientes>"
HEADER = "name: tp0\nservices:\n  server:\n    container_name: server\n    image: server:latest\n    entrypoint: python3 /main.py\n    environment:\n      - PYTHONUNBUFFERED=1\n    networks:\n      - testing_net\n    volumes:\n      - ./server/config.ini:/config.ini:ro\n"
FOOTER = "networks:\n  testing_net:\n    ipam:\n      driver: default\n      config:\n        - subnet: 172.25.125.0/24\n"
DATA = ['    environment:\n      - CLI_ID=1\n      - CLI_LOG_LEVEL=DEBUG\n      - NOMBRE=Rafael\n      - APELLIDO=Putaro\n      - DOCUMENTO=29752645\n      - NACIMIENTO=1982-08-28\n      - NUMERO=7010\n',
        '    environment:\n      - CLI_ID=2\n      - CLI_LOG_LEVEL=DEBUG\n      - NOMBRE=Lionel\n      - APELLIDO=Messi\n      - DOCUMENTO=33888221\n      - NACIMIENTO=1987-06-24\n      - NUMERO=2022\n',
        '    environment:\n      - CLI_ID=3\n      - CLI_LOG_LEVEL=DEBUG\n      - NOMBRE=Pedro\n      - APELLIDO=Pascual\n      - DOCUMENTO=93142191\n      - NACIMIENTO=1970-12-25\n      - NUMERO=3019\n',
        '    environment:\n      - CLI_ID=4\n      - CLI_LOG_LEVEL=DEBUG\n      - NOMBRE=Vito\n      - APELLIDO=Corleone\n      - DOCUMENTO=10522998\n      - NACIMIENTO=1956-07-03\n      - NUMERO=9001\n',
        '    environment:\n      - CLI_ID=5\n      - CLI_LOG_LEVEL=DEBUG\n      - NOMBRE=Beatrix\n      - APELLIDO=Kiddo\n      - DOCUMENTO=31363992\n      - NACIMIENTO=1990-10-17\n      - NUMERO=5189\n']

def do_generate_compose_file(args):
    file = open(args["file"], "w+")
    file.write(HEADER)
    for client_id in range(1, int(args["clients"])+1):
        file.write(generate_client_code(client_id))
    file.write(FOOTER)
    file.close()

def generate_client_code(client_id):
    return f'  client{client_id}:\n    container_name: client{client_id}\n    image: client:latest\n    entrypoint: /client\n{generate_client_environment(client_id)}    networks:\n      - testing_net\n    volumes:\n      - ./client/config.yaml:/config.yaml:ro\n    depends_on:\n      - server\n'

def generate_client_environment(client_id):
    if client_id < 5:
        return DATA[client_id-1]
    else:
        return generate_dummy_data(client_id)

def generate_dummy_data(client_id):
    return f'    environment:\n      - CLI_ID={client_id}\n      - CLI_LOG_LEVEL=DEBUG\n      - NOMBRE=NombreDummy{client_id}\n      - APELLIDO=ApellidoDummy{client_id}\n      - DOCUMENTO={client_id}\n      - NACIMIENTO={1960+client_id}-08-25\n      - NUMERO={client_id}\n'

def generate_compose_file():
    args = parse_args()
    if args:
        do_generate_compose_file(args)

def parse_args():
    if len(sys.argv) >= 3:
        return {
            "file": sys.argv[1],
            "clients": sys.argv[2]
        }
    else:
        print(ERROR_MESSAGE)
        return
    
generate_compose_file()