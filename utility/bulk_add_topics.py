# this program adds topics directly to the database in bulk
# used these resources to help create this program:
#   Cloud SQL docs from Google: https://cloud.google.com/sql/docs/postgres/connect-admin-ip
#   psycopg documentation: https://www.psycopg.org/docs/usage.html 

import csv
import psycopg2
from datetime import datetime, timezone

# connect to the database using psycopg
conn = psycopg2.connect(
    dbname="DB_NAME",
    user="DB_USER",
    password="DB_PASSWORD",
    host="DB_HOST"
)

cur = conn.cursor()

# parse each row and insert into the database
with open("topics.csv", newline='', encoding='utf-8') as csvfile:
    reader = csv.DictReader(csvfile)
    count = 0
    for row in reader:
        cur.execute("""
            INSERT INTO discussion_topics (title, prompt, created_at, is_active, show_count)
            VALUES (%s, %s, %s, %s, %s)
        """, (
            row['title'],
            row['prompt'],
            datetime.now(timezone.utc),
            False,
            0
        ))
        count += 1

conn.commit()
cur.close()
conn.close()