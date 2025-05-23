import psycopg2

conn = psycopg2.connect(
    dbname="capstone",
    user="postgres",
    password="alessandroleatherman",
    host="/tmp/cloudsql/capstone-457718:us-west2:capstone-db"
)

cur = conn.cursor()

cur.execute("DELETE FROM discussion_topics;")

conn.commit()
cur.close()
conn.close()