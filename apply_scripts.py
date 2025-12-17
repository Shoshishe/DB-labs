import psycopg2
import os

conn_string = "host='localhost' dbname='iis' user='shosh' port='5432' password='loshpedus'"
conn = psycopg2.connect(conn_string)

for root, dirs, files in os.walk("sql"):
    for file in files:
        print(os.path.join(root, file))
        conn.cursor().execute(open(os.path.join(root, file), "r").read())
