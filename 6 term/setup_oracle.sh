docker run --name oracle_db -p 1521:1521 -p 5500:5500 \
-e ORACLE_PWD=loshpedus -e ORACLE_CHARACTERSET=UTF8 \
-v oracle-data:/opt/oracle/oradata container-registry.oracle.com/database/express:21.3.0-xe
