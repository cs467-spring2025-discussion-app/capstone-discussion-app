import psycopg2

conn = psycopg2.connect(
    dbname="DB_NAME",
    user="DB_USER",
    password="DB_PASSWORD",
    host="DB_HOST"
)

cur = conn.cursor()

cur.execute("DELETE FROM discussion_topics;")

conn.commit()
cur.close()
conn.close()